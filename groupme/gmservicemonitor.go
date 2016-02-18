package groupme

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/kisielk/sqlstruct"
)

type GroupmeMonitor struct{}

func (g GroupmeMonitor) ValueFor(cachedID int) int {
	posts, _ := postFetch("WHERE cache_id = $1", []interface{}{cachedID})

	groupedPosts := make(map[string][]string)
	for _, post := range posts {
		slice := groupedPosts[post.GroupID]
		slice = append(slice, post.MessageID)
		groupedPosts[post.GroupID] = slice
	}

	token := os.Getenv("groupMeID")

	var wg sync.WaitGroup
	out := make(chan int)

	for key, ids := range groupedPosts {
		wg.Add(1)
		go func(groupID string, messageIDs []string) {
			getURL := "https://api.groupme.com/v3/groups/" + groupID + "/likes?period=day&token=" + token
			resp, _ := http.Get(getURL)
			body, _ := ioutil.ReadAll(resp.Body)

			var wrapper gmMessageWrapper
			json.Unmarshal(body, &wrapper)

			for _, message := range wrapper.Response.Messages {
				for _, id := range messageIDs {
					if id == message.MessageID() {
						out <- len(message.FavoritedBy)
					}
				}
			}
			wg.Done()
		}(key, ids)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	likes := 0
	for result := range out {
		likes += result
	}

	fmt.Print("LIKES - ")
	fmt.Println(likes)

	return likes
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

//Temp DB
var currentDB *sql.DB

func init() {
	currentDB = connect()
}

func connect() *sql.DB {
	dbUrl := os.Getenv("DATABASE_URL")
	database, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("[x] Could not open the connection to the database. Reason: %s", err.Error())
	}
	return database
}

func postFetch(whereStr string, values []interface{}) ([]GroupmePost, error) {

	queryStr := fmt.Sprintf("SELECT %s FROM groupme_posts", sqlstruct.Columns(GroupmePost{}))

	fmt.Println(queryStr)

	rows, err := currentDB.Query(queryStr+" "+whereStr, values...)
	if err != nil {
		return []GroupmePost{}, err
	}
	defer rows.Close()

	actions := []GroupmePost{}
	for rows.Next() {
		var act GroupmePost
		err := sqlstruct.Scan(&act, rows)
		if err == nil {
			actions = append(actions, act)
		}
	}

	return actions, nil
}
