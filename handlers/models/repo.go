package models

import (
    "database/sql"
    "os"
    "log"
    
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

func Name() string {
    var name string
    currentDB.QueryRow("SELECT name from bots WHERE id = '1'").Scan(&name)
    
    return name
}

func FetchBot(id string) (Bot, error) {
    row := currentDB.QueryRow("SELECT * from bots WHERE group_id = $1", id)
    
    var bot Bot
    err := row.Scan(&bot.GroupID, &bot.GroupName, &bot.BotName, &bot.Key)
    if err != nil {
        return Bot{}, err
    }
    
    return bot, nil
}

func FetchActions(primary bool) ([]Action, error) {
    
    queryStr := "SELECT type, content, data_path, pattern, main, priority, fallback_action from actions"
    if primary {
        queryStr += " WHERE main = 'TRUE'"
    }
    rows, err := currentDB.Query(queryStr)
    defer rows.Close()
    
    if err != nil {
        return make([]Action, 0), err
    }
    
    actions := make([]Action, 1)
    for rows.Next() {
        var act Action
        err := rows.Scan(&act.ContentType, &act.Content, &act.DataPath, &act.Pattern, &act.Primary, &act.Priority, &act.FallbackAction)
        if err != nil {
            continue
        }
        
        actions = append(actions, act)
    }
    
    return actions, nil
}