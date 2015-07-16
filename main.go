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
    
    var server http.Handler
    path := r.URL.Path
    if path == "/addBot" {
        server = ABHook{}
    } else {
        server = GMHook{} 
    }
    
    server.ServeHTTP(w, r)
}

func main() {
    var l LocalHandler
    err := http.ListenAndServe(GetPort(), l)
	if err != nil {
		log.Fatal(err)
	}
}

func GetPort() string {
    port := os.Getenv("PORT")
    if port == "" {
        port = "4000"
    }
    
    return ":" + port
}