package groupme

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"
)

type gmPost struct {
	ID        int    `sql:"id"`
	CacheID   int    `sql:"cache_id"`
	Likes     int    `sql:"likes"`
	MessageID string `sql:"message_id"`
	GroupID   string `sql:"group_id"`
}

func cachePost(cacheID int, messageID, groupID string) {
	queryStr := "INSERT INTO groupme_posts(cache_id, message_id, group_id) VALUES($1, $2, $3)"
	config.DB.QueryRow(queryStr, cacheID, messageID, groupID)
}

func updateLikes() {
	queryStr := fmt.Sprintf("SELECT %s FROM groupme_posts WHERE posted_at >= NOW() - '1 day'::INTERVAL", sqlstruct.Columns(gmPost{}))

	rows, err := config.DB.Query(queryStr)
	if err != nil {
		return
	}
	defer rows.Close()

	groupedPosts := make(map[string][]gmPost)
	for rows.Next() {
		var post gmPost
		err := sqlstruct.Scan(&post, rows)
		if err == nil {
			slice := groupedPosts[post.GroupID]
			slice = append(slice, post)
			groupedPosts[post.GroupID] = slice
		}
	}

	token := os.Getenv("groupMeID")

	updated := make(map[int]int)

	for key, group := range groupedPosts {
		getURL := "https://api.groupme.com/v3/groups/" + key + "/likes?period=day&token=" + token
		resp, _ := http.Get(getURL)
		body, _ := ioutil.ReadAll(resp.Body)

		var wrapper gmMessageWrapper
		json.Unmarshal(body, &wrapper)

		for _, message := range wrapper.Response.Messages {
			for _, post := range group {
				if post.MessageID == message.MessageID() {
					updated[post.ID] = len(message.FavoritedBy)
				}
			}
		}
	}

	tx, err := config.DB.Begin()
	if err != nil {
		return
	}

	stmt, _ := config.DB.Prepare("UPDATE groupme_posts SET likes=$1 WHERE id=$2")
	for updateID, likeCount := range updated {
		stmt.Exec(likeCount, updateID)
	}
	tx.Commit()
}
