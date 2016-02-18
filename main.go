package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	router := NewRouter()

	log.Fatal(http.ListenAndServe(GetPort(), router))
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	return ":" + port
}
