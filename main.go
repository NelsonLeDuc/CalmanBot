package main

import (
	"log"
	"net/http"
    "os"
)

func main() {
	var g GMHook
	err := http.ListenAndServe(":"+os.Getenv("PORT"), g)
	if err != nil {
		log.Fatal(err)
	}
}