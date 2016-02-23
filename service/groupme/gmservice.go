package groupme

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/nelsonleduc/calmanbot/service"
	"github.com/nelsonleduc/calmanbot/utility"
)

const postDelayMilliseconds = 500
const messagePollDelayMilliseconds = 500
const groupmeLengthLimit = 1000

func init() {
	service.AddService("groupme", gmService{})
}

type gmService struct{}

func (g gmService) PostText(key, text string, cacheID int, groupMessage service.Message) {

	dividedText := utility.DivideString(text, groupmeLengthLimit)
	if len(dividedText) > 2 {
		pasteBinURL := postToPastebin(text)
		if len(pasteBinURL) == 0 {
			return
		}

		dividedText = []string{pasteBinURL}
	}

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

			postToGroupMe(encoded)
			mID, err := messageID(groupMessage)
			if err == nil {
				cachePost(cacheID, mID, groupMessage.GroupID())
			}
		}(key, subText)
	}

	go updateLikes()
}

func (g gmService) MessageFromJSON(reader io.Reader) service.Message {
	message := new(gmMessage)
	json.NewDecoder(reader).Decode(message)

	return *message
}

func (g gmService) ServiceMonitor() (service.Monitor, error) {
	return GroupmeMonitor{}, nil
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

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return string(respBody)
}

func postToGroupMe(body []byte) {
	time.Sleep(postDelayMilliseconds * time.Millisecond)

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
	getURL := "https://api.groupme.com/v3/groups/" + message.GroupID() + "/messages?token=" + token + "&after_id=" + message.MessageID()
	resp, _ := http.Get(getURL)

	body, _ := ioutil.ReadAll(resp.Body)

	var wrapper gmMessageWrapper
	json.Unmarshal(body, &wrapper)

	for _, recieved := range wrapper.Response.Messages {
		if recieved.UserType() == "bot" {
			return recieved.MessageID(), nil
		}
	}

	return "", errors.New("No messages found")
}
