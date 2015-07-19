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
    "regexp"
    "sort"
)

func isValidHTTPURLString(s string)  bool {
    URL, _ := url.Parse(s)
    return (URL.Scheme == "http" || URL.Scheme == "https")
}

func HandleCalman(w http.ResponseWriter, r *http.Request) {
    
    message := ParseMessageJSON(r.Body)
    bot, _ := models.FetchBot(message.GroupID)
    
    if !strings.HasPrefix(message.Text, "@" + bot.BotName) {
        return
    }
    
    actions, _ := models.FetchActions(true)
    sort.Sort(models.ByPriority(actions))
    
    var (
        act models.Action
        sMatch string
    )
    for _, a := range actions {
        r, _ := regexp.Compile(*a.Pattern)
        matched := r.FindStringSubmatch(message.Text)
        if len(matched) > 1 && matched[1] != "" {
            sMatch = matched[1]
            act = a
            break
        }
    }
    
    updateAction(&act, sMatch)
    
    if act.IsURLType() {
        handleURLAction(act, w, bot)
    }
}

func handleURLAction(a models.Action, w http.ResponseWriter, b models.Bot) {
    
    fmt.Fprintln(w, a)
    resp, err := http.Get(a.Content)
    
    failure := func() {
        
        fmt.Printf("Failed")
    }
    
    if err == nil {
        
        content, _ := ioutil.ReadAll(resp.Body)
        pathString := *a.DataPath
        
        str := ParseJSON(content, pathString)
        if str == "" {
            failure()
        }
        
        success := func(s string) {
            fmt.Printf("Success: %v\n", s)
            postText(b, s)
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
    } else {
        failure()
    }
    
    resp.Body.Close()
}

func postText(b models.Bot, t string) {
    
    postURL := "https://api.groupme.com/v3/bots/post"
    postBody := map[string]string {
        "bot_id": b.Key,
        "text": t,
    }
    
    encoded, _ := json.Marshal(postBody)
    
    http.Post(postURL, "application/json", bytes.NewReader(encoded))
}

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

func updateAction(a *models.Action, text string) {
    text = url.QueryEscape(text)
    
    a.Content = strings.Replace(a.Content, "{_text_}", text, -1)
}