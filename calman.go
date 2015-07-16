package main

import (
	"fmt"
	"net/http"
    "net/url"
//    "strings"
    "io/ioutil"
)

type GMHook struct {}

type Message struct {
    GroupID string      `json:"group_id"`
    UserName string     `json:"name"`
    UserID string       `json:"id"`
    Text string         `json:"text"`
    UserType string     `json:"sender_type"`
}

func isValidHTTPURLString(s string)  bool {
    URL, _ := url.Parse(s)
    return (URL.Scheme == "http" || URL.Scheme == "https")
}


func (gm GMHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//    bytes := []byte(`{"attachments": [], "avatar_url": "http://i.groupme.com/123456789", "created_at": 1302623328, "group_id": "1234567890", "id": "1234567890", "name": "John", "sender_id": "12345", "sender_type": "user", "source_guid": "GUID", "system": false, "text": "Hello world ☃☃", "user_id": "1234567890"}`)
    
//    reader := strings.NewReader(string(bytes))
//    message := ParseMessageJSON(reader)
    
    resp, err := http.Get("http://ajax.googleapis.com/ajax/services/search/images?v=1.0&as_filetype=gif&imgtype=animated&rsz=8&q=ambiguity%20strikes%20again")
    if err == nil {
        
        content, _ := ioutil.ReadAll(resp.Body)
        pathString := "responseData.results.{_randomInt_}.url"
        
        str := ParseJSON(content, pathString)
        
        success := func(s string) {
            fmt.Printf("Success: %v\n", s)
            fmt.Fprintln(w, s)
        }
        failure := func() {
            //Actually perform fallback here
            
            fmt.Printf("Failed")
        }
        
        if !ValidateURL(str, success) {
            fmt.Printf("Invalid URL: %v\n", str)
            
            oldStr := str
            for i := 0; i < 3 && oldStr == str; i++ {
                str = ParseJSON(content, pathString)
            }
            
            if !ValidateURL(str, success) {
                failure()
            }
        }
    }
    
    resp.Body.Close()
}


func ValidateURL(u string, success func(string)) bool {
    if isValidHTTPURLString(u) {
        resp, err := http.Get(u)
        
        if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
            success(u)
        } else {
            return false
        }
    } else {
        success(u)
    }
    
    return true
}