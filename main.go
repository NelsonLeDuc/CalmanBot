package main

import (
	"log"
	"net/http"
)

func main() {
	var g GMHook
	err := http.ListenAndServe("localhost:4000", g)
	if err != nil {
		log.Fatal(err)
	}
}