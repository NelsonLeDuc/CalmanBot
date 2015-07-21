package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

func isValidHTTPURLString(s string) bool {
	URL, _ := url.Parse(s)
	return (URL.Scheme == "http" || URL.Scheme == "https")
}

func HandleCalman(w http.ResponseWriter, r *http.Request) {

	message := ParseMessageJSON(r.Body)
	bot, _ := models.FetchBot(message.GroupID)

	if len(message.Text) < 1 || !strings.HasPrefix(strings.ToLower(message.Text[1:]), strings.ToLower(bot.BotName)) {
		return
	}

	actions, _ := models.FetchActions(true)
	sort.Sort(models.ByPriority(actions))

	var (
		act    models.Action
		sMatch string
	)
	for _, a := range actions {
		r, _ := regexp.Compile("(?i)" + *a.Pattern)
		matched := r.FindStringSubmatch(message.Text)
		if len(matched) > 1 && matched[1] != "" {
			sMatch = matched[1]
			act = a
			break
		}
	}

	postString := ""
	for {
		updateAction(&act, sMatch)
		if act.IsURLType() {
			postString = handleURLAction(act, w, bot)
		} else {
			postString = act.Content
		}

		if postString != "" || act.FallbackAction == nil {
			break
		} else {
			act, _ = models.FetchAction(*act.FallbackAction)
		}
	}

	if postString != "" {
		fmt.Printf("Action: %v\n", act.Content)
		fmt.Printf("Posting: %v\n", postString)
		postText(bot, postString)
	}
}

func handleURLAction(a models.Action, w http.ResponseWriter, b models.Bot) string {

	fmt.Fprintln(w, a)
	resp, err := http.Get(a.Content)

	if err == nil {

		content, _ := ioutil.ReadAll(resp.Body)
		pathString := *a.DataPath

		str := ParseJSON(content, pathString)
		if str == "" {
			return ""
		} else {

			if !validateURL(str, a.IsImageType()) {
				fmt.Printf("Invalid URL: %v\n", str)

				oldStr := str
				for i := 0; i < 3 && oldStr == str; i++ {
					str = ParseJSON(content, pathString)
				}

				if !validateURL(str, a.IsImageType()) {
					return ""
				} else {
					return str
				}
			} else {
				return str
			}
		}
	} else {
		return ""
	}

	resp.Body.Close()
	return ""
}

func postText(b models.Bot, t string) {

	postURL := "https://api.groupme.com/v3/bots/post"
	postBody := map[string]string{
		"bot_id": b.Key,
		"text":   t,
	}

	encoded, _ := json.Marshal(postBody)

	http.Post(postURL, "application/json", bytes.NewReader(encoded))
}

func validateURL(u string, image bool) bool {

	if isValidHTTPURLString(u) {

		resp, err := http.Get(u)
		defer resp.Body.Close()

		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if image {
				return validateImage(resp.Body)
			}

			return true
		} else {
			return false
		}
	} else {
		return true
	}

	return true
}

func updateAction(a *models.Action, text string) {
	text = url.QueryEscape(text)

	a.Content = strings.Replace(a.Content, "{_text_}", text, -1)
}

//TODO: Move this out of here
func validateImage(r io.Reader) bool {
	buf := make([]byte, 8)
	num, err := r.Read(buf)
	if err != nil || num < 8 {
		return false
	}

	gif := convertHexSlice([]string{"47", "49", "46"})
	jpg := convertHexSlice([]string{"FF", "D8", "FF", "E0"})
	png := convertHexSlice([]string{"89", "50", "4E", "47", "D", "A", "1A", "A"})

	return byteSliceSubset(buf, gif) || byteSliceSubset(buf, jpg) || byteSliceSubset(buf, png)
}

//b is a subset of a
func byteSliceSubset(a, b []byte) bool {
	for i, el := range b {
		if el != a[i] {
			return false
		}
	}
	
	fmt.Printf("Image valid, bytes: %v\n", b)
	return true
}

func convertHexSlice(s []string) []byte {
	b := []byte{}
	for _, hexStr := range s {
		result, err := hex.DecodeString(hexStr)
		if err == nil {
			b = append(b, result[0])
		}
	}

	return b
}
