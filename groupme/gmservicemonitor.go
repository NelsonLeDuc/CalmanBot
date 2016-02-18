package groupme

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/kisielk/sqlstruct"
)

type GroupmeMonitor struct{}

func (g GroupmeMonitor) ValueFor(cachedID int) int {
	row := currentDB.QueryRow("SELECT sum(likes) FROM groupme_posts WHERE cache_id=$1 GROUP BY cache_id", cachedID)

	var likeCount int
	row.Scan(&likeCount)

	return likeCount
}

func cachePost(cacheID int, messageID, groupID string) {
	queryStr := "INSERT INTO groupme_posts(cache_id, message_id, group_id) VALUES($1, $2, $3)"
	currentDB.QueryRow(queryStr, cacheID, messageID, groupID)
}

type GroupmePost struct {
	ID        int    `sql:"id"`
	CacheID   int    `sql:"cache_id"`
	Likes     int    `sql:"likes"`
	MessageID string `sql:"message_id"`
	GroupID   string `sql:"group_id"`
}

func updateLikes() {
	queryStr := fmt.Sprintf("SELECT %s FROM groupme_posts WHERE posted_at >= NOW() - '1 day'::INTERVAL", sqlstruct.Columns(GroupmePost{}))

	rows, err := currentDB.Query(queryStr)
	if err != nil {
		return
	}
	defer rows.Close()

	groupedPosts := make(map[string][]GroupmePost)
	for rows.Next() {
		var post GroupmePost
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

	tx, err := currentDB.Begin()
	if err != nil {
		return
	}

	stmt, _ := currentDB.Prepare("UPDATE groupme_posts SET likes=$1 WHERE id=$2")
	for updateID, likeCount := range updated {
		stmt.Exec(likeCount, updateID)
	}
	tx.Commit()
}

//Temp DB
var currentDB *sql.DB

func init() {
	dbUrl := os.Getenv("DATABASE_URL")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("[x] Could not open the connection to the database. Reason: %s", err.Error())
	}

	currentDB = database
}
