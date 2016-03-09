package models

import (
	"fmt"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"
)

func actionFetch(whereStr string, values []interface{}) ([]Action, error) {

	queryStr := fmt.Sprintf("SELECT %s FROM actions", sqlstruct.Columns(Action{}))

	rows, err := config.DB.Query(queryStr+" "+whereStr, values...)
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
	rows, err := config.DB.Query(fmt.Sprintf("SELECT %s FROM bots WHERE group_id = $1", sqlstruct.Columns(Bot{})), id)
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
