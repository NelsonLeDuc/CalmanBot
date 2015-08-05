package main

import (
	"github.com/nelsonleduc/calmanbot/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"CalmanRespond",
		"POST",
		"/bots/botHook",
		handlers.HandleBotHook,
	},
	Route{
		"ActionsRespond",
		"GET",
		"/actions",
		handlers.HandleActions,
	},
}
