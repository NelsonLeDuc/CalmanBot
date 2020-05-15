package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/handlers"

	"github.com/gorilla/mux"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	router := NewRouter()

	if config.Configuration().EnableDiscord() {
		go CreateWebhook()
	}

	if config.Configuration().EnableMinecraft() {
		go handlers.MonitorMinecraft()
	}

	log.Fatal(http.ListenAndServe(GetPort(), router))
}

func GetPort() string {
	port := config.Configuration().Port()

	if config.Configuration().VerboseMode() {
		fmt.Printf("Listening on port %v\n\n", port)
	}

	return ":" + port
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	verbose := config.Configuration().VerboseMode()
	if verbose {
		fmt.Println("=== Setup Router ===")
	}
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

		if verbose {
			fmt.Println("      Name: ", route.Name)
			fmt.Println("      Path: ", route.Pattern)
			fmt.Printf("    Method:  %v\n\n", route.Method)
		}
	}

	//Serve static content from the static directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return router
}
