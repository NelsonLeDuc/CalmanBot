package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/nelsonleduc/calmanbot/service"
	"github.com/nelsonleduc/calmanbot/utility"
)

const postDelayMilliseconds = 500
const groupmeLengthLimit = 1000

func init() {
	service.AddService("groupme", gmService{})
}

type gmService struct{}

func (g gmService) PostText(key, text string, cacheID int, groupMessage service.Message) {

	dividedText := utility.DivideString(text, groupmeLengthLimit)

	for _, subText := range dividedText {
		func(key, message string) {
			postBody := map[string]string{
				"bot_id": key,
				"text":   message,
			}

			encoded, err := json.Marshal(postBody)
			if err != nil {
				return
			}

			postToGroupMe(encoded)
			cachePost(cacheID, key)
			messageID(groupMessage)
		}(key, subText)
	}
}

type gmMessageWrapper struct {
	Response struct {
		Count    int         `json:"count"`
		Messages []gmMessage `json:"messages"`
	} `json:"response"`
}

func messageID(message service.Message) string {
	token := os.Getenv("groupMeID")
	getURL := "https://api.groupme.com/v3/groups/" + message.GroupID() + "/messages?token=" + token + "&after_id=" + message.MessageID()
	resp, _ := http.Get(getURL)

	body, _ := ioutil.ReadAll(resp.Body)

	var wrapper gmMessageWrapper
	json.Unmarshal(body, &wrapper)

	fmt.Println(body)

	fmt.Println(getURL)
	fmt.Println("messages:")
	fmt.Println(wrapper.Response.Messages)

	return ""
}

func (g gmService) MessageFromJSON(reader io.Reader) service.Message {
	message := new(gmMessage)
	json.NewDecoder(reader).Decode(message)

	return *message
}

func (g gmService) ServiceMonitor() *service.Monitor {
	return nil
}

func postToGroupMe(body []byte) {
	time.Sleep(postDelayMilliseconds * time.Millisecond)

	postURL := "https://api.groupme.com/v3/bots/post"
	http.Post(postURL, "application/json", bytes.NewReader(body))
}
