package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func HandleYoutubeLinkt(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("link")
	if query == "" {
		return
	}

	url := "https://api.song.link/v1-alpha.1/links?url=" + url.QueryEscape(query)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "CalmanBot/2.5.3")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Bad Response: %d length: %d", resp.StatusCode, len(content))
		return
	}

	var stuff map[string]interface{}

	err = json.Unmarshal(content, &stuff)
	if err != nil {
		return
	}
	pageURL := stuff["pageUrl"].(string)
	entities := stuff["entitiesByUniqueId"].(map[string]interface{})
	if len(entities) <= 2 {
		return
	}
	outputData := map[string]string{
		"pageUrl": pageURL,
	}

	json, _ := json.Marshal(outputData)
	w.Write(json)
}
