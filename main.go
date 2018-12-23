package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nelsonleduc/calmanbot/config"

	"github.com/gorilla/mux"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	router := NewRouter()

	if config.Configuration().EnableDiscord() {
		CreateWebhook()
	}

	log.Fatal(http.ListenAndServe(GetPort(), router))
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	return ":" + port
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	//Serve static content from the static directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router
}
