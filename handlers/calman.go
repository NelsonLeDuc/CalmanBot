package handlers

import (
	"fmt"
	"net/http"
    "net/url"
    "io/ioutil"
    "strings"
    "github.com/nelsonleduc/calmanbot/handlers/models"
)

type GMHook struct {}

func isValidHTTPURLString(s string)  bool {
    URL, _ := url.Parse(s)
    return (URL.Scheme == "http" || URL.Scheme == "https")
}

func HandleCalman(w http.ResponseWriter, r *http.Request) {
    
    act, _ := models.FetchAction(12)
    UpdateAction(&act, models.Message{Text: "cats"})
    
    if act.IsURLType() {
        HandleURLAction(act, w)
    }
}

func PrintMessage(w http.ResponseWriter, r *http.Request) {
    
    cont, _ := ioutil.ReadAll(r.Body)
    fmt.Println(string(cont))
}

func HandleURLAction(a models.Action, w http.ResponseWriter) {
    
    fmt.Fprintln(w, a)
    resp, err := http.Get(a.Content)
    if err == nil {
        
        content, _ := ioutil.ReadAll(resp.Body)
        pathString := *a.DataPath
        
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
    
    client := http.Client{}
    if isValidHTTPURLString(u) {
        req, err := http.NewRequest("HEAD", u, nil)
        if err != nil {
            return false
        }
        
        resp, err := client.Do(req)
        
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

func UpdateAction(a *models.Action, m models.Message) {
    text := url.QueryEscape(m.Text)
    
    a.Content = strings.Replace(a.Content, "{_text_}", text, -1)
}