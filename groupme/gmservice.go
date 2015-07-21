package groupme

import (
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	
	"github.com/nelsonleduc/calmanbot/service"
)

func init() {
	service.AddService("groupme", gmService{})
}

type gmService struct {}

func (g gmService) PostText(key, text string) error {

	postURL := "https://api.groupme.com/v3/bots/post"
	postBody := map[string]string{
		"bot_id": key,
		"text":   text,
	}

	encoded, err := json.Marshal(postBody)
	if err != nil {
		return err
	}
	
	_, err = http.Post(postURL, "application/json", bytes.NewReader(encoded))
	
	return err
}

func (g gmService) MessageFromJSON(reader io.Reader) service.Message {
	message := new(gmMessage)
	json.NewDecoder(reader).Decode(message)
	
	return *message
}
