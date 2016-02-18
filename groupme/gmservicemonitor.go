package groupme

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/kisielk/sqlstruct"
)

type GroupmeMonitor struct{}

func (g GroupmeMonitor) ValueFor(cachedID int) int {
	return 1
}

func cachePost(cacheID int, messageID string) {
	queryStr := "INSERT INTO groupme_posts(cache_id, message_id) VALUES($1, $2)"
	currentDB.QueryRow(queryStr, cacheID, messageID)
}

type GroupmePost struct {
	ID      int    `sql:"id"`
	CacheID int    `sql:"cache_id"`
	Key     string `sql:"key"`
	Likes   int    `sql:"likes"`
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

	queryStr := fmt.Sprintf("SELECT %s FROM cached", sqlstruct.Columns(GroupmePost{}))

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
