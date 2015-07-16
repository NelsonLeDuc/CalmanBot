package main

import (
//	"fmt"
	"net/http"
//    "net/url"
//    "strings"
//    "io/ioutil"
)

type ABHook struct {}

func (ab ABHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    
    if r.Method == "GET" {
        http.ServeFile(w, r, "html/addBot.html")
    } else if r.Method == "POST" {
        
    }
}