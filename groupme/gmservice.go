package groupme

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	"time"
	
	"github.com/nelsonleduc/calmanbot/service"
)

const postDelayMilliseconds = 500

func init() {
	service.AddService("groupme", gmService{})
}

type gmService struct {}

func (g gmService) PostText(key, text string) {
	postBody := map[string]string{
		"bot_id": key,
		"text":   text,
	}

	encoded, err := json.Marshal(postBody)
	if err != nil {
		return
	}
	
	go postToGroupMe(encoded)
}

func (g gmService) MessageFromJSON(reader io.Reader) service.Message {
	message := new(gmMessage)
	json.NewDecoder(reader).Decode(message)
	
	return *message
}

func postToGroupMe(body []byte) {
	time.Sleep(postDelayMilliseconds * time.Millisecond)

	postURL := "https://api.groupme.com/v3/bots/post"
	http.Post(postURL, "application/json", bytes.NewReader(body))
}
