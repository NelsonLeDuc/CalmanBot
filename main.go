package main

import (
	"log"
	"net/http"
    "os"
    "database/sql"
)

func main() {
	var g GMHook
    err := http.ListenAndServe(GetPort(), g)
	if err != nil {
		log.Fatal(err)
	}
    
    dbUrl := os.Getenv("DATABASE_URL")
    database, _ := sql.Open("postgres", dbUrl)
    log.Fatal(database)
}

func GetPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }
    
    return ":" + port
}