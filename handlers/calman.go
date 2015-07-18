package handlers

import (
	"fmt"
	"net/http"
    "net/url"
    "io/ioutil"
    "strings"
    "github.com/nelsonleduc/calmanbot/handlers/models"
    "encoding/json"
    "bytes"
)

type GMHook struct {}

func isValidHTTPURLString(s string)  bool {
    URL, _ := url.Parse(s)
    return (URL.Scheme == "http" || URL.Scheme == "https")
}

func HandleCalman(w http.ResponseWriter, r *http.Request) {
    
    message := ParseMessageJSON(r.Body)
    
    act, _ := models.FetchAction(12)
    updateAction(&act, message)
    
    bot, _ := models.FetchBot(message.GroupID)
    
    if act.IsURLType() {
        handleURLAction(act, w, bot)
    }
}

func handleURLAction(a models.Action, w http.ResponseWriter, b models.Bot) {
    
    fmt.Fprintln(w, a)
    resp, err := http.Get(a.Content)
    if err == nil {
        
        content, _ := ioutil.ReadAll(resp.Body)
        pathString := *a.DataPath
        
        str := ParseJSON(content, pathString)
        
        success := func(s string) {
            fmt.Printf("Success: %v\n", s)
//            fmt.Fprintln(w, s)
            postText(b, s)
        }
        failure := func() {
            //Actually perform fallback here
            
            fmt.Printf("Failed")
        }
        
        if !validateURL(str, success) {
            fmt.Printf("Invalid URL: %v\n", str)
            
            oldStr := str
            for i := 0; i < 3 && oldStr == str; i++ {
                str = ParseJSON(content, pathString)
            }
            
            if !validateURL(str, success) {
                failure()
            }
        }
    }
    
    resp.Body.Close()
}

func postText(b models.Bot, t string) {
    
    t = url.QueryEscape(t)
    postURL := "https://api.groupme.com/v3/bots/post"
    postBody := map[string]string {
        "bot_id": b.Key,
        "text": t,
    }
    
    encoded, _ := json.Marshal(postBody)
    
    http.Post(postURL, "application/json", bytes.NewReader(encoded))
}

//    Parse.Cloud.httpRequest({
//        url: "https://api.groupme.com/v3/bots/post?bot_id=" + gBot.key + "&text=" + encodeURIComponent(text),
//        method: "POST",
//        success: function (httpResponse) {
//            var GroupMessage = Parse.Object.extend("GroupMessage");
//            var groupMessage = new GroupMessage();
// 
//            groupMessage.save({
//                text: original,
//                user: gUser.name,
//                imageURL: text,
//                groupIdentifier: gBot.groupID,
//                userIdentifier: gUser.ID
//            });
//            res.send("Done Posting Image")
//        },
//        error: function (httpResponse) {
//            res.send(418, "Stop brewing me!")
//        }
//    });

func validateURL(u string, success func(string)) bool {
    
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

func updateAction(a *models.Action, m models.Message) {
    text := url.QueryEscape(m.Text)
    
    a.Content = strings.Replace(a.Content, "{_text_}", text, -1)
}