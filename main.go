package main

import (
	"log"
	"net/http"
    "os"
//    "database/sql"
//    "fmt"
)

type LocalHandler struct {}

func (l LocalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "html/addBot.html")
}

func main() {
	var g GMHook
    var l LocalHandler
    
    err := http.ListenAndServe(GetPort(), g)
	if err != nil {
		log.Fatal(err)
	}
    
    http.ListenAndServe("/addBot" + GetPort(), l)
    
//    dbUrl := os.Getenv("DATABASE_URL")
//    database, _ := sql.Open("postgres", dbUrl)
}

func GetPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }
    
    return ":" + port
}