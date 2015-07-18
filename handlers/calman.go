package handlers

import (
	"fmt"
	"net/http"
    "net/url"
    "io/ioutil"
    "github.com/nelsonleduc/calmanbot/handlers/models"
)

type GMHook struct {}

func isValidHTTPURLString(s string)  bool {
    URL, _ := url.Parse(s)
    return (URL.Scheme == "http" || URL.Scheme == "https")
}

func HandleCalman(w http.ResponseWriter, r *http.Request) {
    
    resp, err := http.Get("http://ajax.googleapis.com/ajax/services/search/images?v=1.0&as_filetype=gif&imgtype=animated&rsz=8&q=ambiguity%20strikes%20again")
    if err == nil {
        
        bot, _ := models.FetchActions(true)//.FetchBot("9214876")
        fmt.Fprintln(w, bot)
        
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