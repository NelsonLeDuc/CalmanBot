package main

import (
    "net/http"
    "github.com/nelsonleduc/calmanbot/handlers"
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
        "/botHook",
        handlers.HandleCalman,
    },
}