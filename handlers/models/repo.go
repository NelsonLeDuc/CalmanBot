package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/kisielk/sqlstruct"
	_ "github.com/lib/pq"
)

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

func actionFetch(whereStr string, values []interface{}) ([]Action, error) {

	queryStr := fmt.Sprintf("SELECT %s FROM actions", sqlstruct.Columns(Action{}))

	rows, err := currentDB.Query(queryStr+" "+whereStr, values...)
	if err != nil {
		return []Action{}, err
	}
	defer rows.Close()

	actions := []Action{}
	for rows.Next() {
		var act Action
		err := sqlstruct.Scan(&act, rows)
		if err == nil {
			actions = append(actions, act)
		}
	}

	return actions, nil
}

//Public Methods
func FetchBot(id string) (Bot, error) {
	rows, err := currentDB.Query(fmt.Sprintf("SELECT %s FROM bots WHERE group_id = $1", sqlstruct.Columns(Bot{})), id)
	if err != nil {
		return Bot{}, err
	}
	defer rows.Close()

	rows.Next()
	var bot Bot
	err = sqlstruct.Scan(&bot, rows)

	return bot, err
}

func FetchActions(primary bool) ([]Action, error) {
	var (
		values   []interface{}
		queryStr string
	)
	if primary {
		queryStr = "WHERE main = $1"
		values = append(values, primary)
	}

	return actionFetch(queryStr, values)
}

func FetchAction(id int) (Action, error) {
	actions, err := actionFetch("WHERE id = $1", []interface{}{id})

	var action Action
	if len(actions) > 0 {
		action = actions[0]
	}

	return action, err
}
