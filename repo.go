package main

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