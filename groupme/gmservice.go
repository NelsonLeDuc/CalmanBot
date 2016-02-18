package groupme

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

func (g gmService) PostText(key, text string, cacheID int) {

	dividedText := utility.DivideString(text, groupmeLengthLimit)

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
		}(key, subText)
	}
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
