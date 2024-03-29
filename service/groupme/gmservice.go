package groupme

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
	"github.com/nelsonleduc/calmanbot/utility"
)

const postDelayMilliseconds = 500
const messagePollDelayMilliseconds = 500
const groupmeLengthLimit = 1000

type GMService struct{}

func (g GMService) Post(post service.Post, groupMessage service.Message) { //PostText(key, text string, pType service.PostType, cacheID int, groupMessage service.Message) {

	if config.Configuration().VerboseMode() {
		log.Print("***POSTING***")
		log.Print("key:  " + post.Key)
		log.Print("text: " + post.Text)
		log.Print("***DONE***")
	}

	dividedText := utility.DivideString(post.Text, groupmeLengthLimit)
	if len(dividedText) > 2 {
		pasteBinURL := postToPastebin(post.Text)
		if len(pasteBinURL) == 0 {
			return
		}

		dividedText = []string{pasteBinURL}
	}

	idx := time.Duration(1)
	for _, subText := range dividedText {
		go func(key, message string) {
			postBody := map[string]string{
				"bot_id": key,
				"text":   message,
			}

			encoded, err := json.Marshal(postBody)
			if err != nil {
				return
			}

			postToGroupMe(encoded, idx)
			mID, err := messageID(groupMessage)
			if err == nil {
				cachePost(post.CacheID, mID, groupMessage.BotGroupID())
			}
		}(post.Key, subText)
		idx++
	}

	go updateLikes()
}

func (g GMService) MessageFromJSON(reader io.Reader) service.Message {
	message := new(gmMessage)
	json.NewDecoder(reader).Decode(message)

	return *message
}

func (g GMService) ServiceMonitor() (service.Monitor, error) {
	return GroupmeMonitor{}, nil
}

func (g GMService) NoteProcessing(groupMessage service.Message) {
	// no-op
}

func (g GMService) ServiceTriggerWrangler() (service.TriggerWrangler, error) {
	return nil, errors.New("Unsupported")
}

func (g GMService) SupportsBuiltInFeature(feature service.BuiltInFeature) bool {
	switch feature {
	case service.BuiltInFeatureLeaderboard:
		return true
	}
	return false
}

func postToPastebin(text string) string {
	data := url.Values{}

	data.Set("api_dev_key", os.Getenv("pastebinKey"))
	data.Set("api_option", "paste")
	data.Set("api_paste_code", text)

	data.Set("api_paste_name", "Too long for /r/shitjustinsays")
	data.Set("api_paste_private", "0")
	data.Set("api_paste_expire_date", "1D")

	resp, err := http.PostForm("http://pastebin.com/api/api_post.php", data)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(respBody)
}

func postToGroupMe(body []byte, multiplier time.Duration) {
	time.Sleep((postDelayMilliseconds * multiplier) * time.Millisecond)

	postURL := "https://api.groupme.com/v3/bots/post"
	http.Post(postURL, "application/json", bytes.NewReader(body))
}

type gmMessageWrapper struct {
	Response struct {
		Messages []gmMessage `json:"messages"`
	} `json:"response"`
}

func messageID(message service.Message) (string, error) {
	time.Sleep(messagePollDelayMilliseconds * time.Millisecond)
	token := os.Getenv("groupMeID")
	getURL := "https://api.groupme.com/v3/groups/" + message.BotGroupID() + "/messages?token=" + token + "&after_id=" + message.MessageID()
	resp, _ := http.Get(getURL)

	body, _ := io.ReadAll(resp.Body)

	var wrapper gmMessageWrapper
	json.Unmarshal(body, &wrapper)

	for _, recieved := range wrapper.Response.Messages {
		if recieved.UserType() == "bot" {
			return recieved.MessageID(), nil
		}
	}

	return "", errors.New("no messages found")
}
