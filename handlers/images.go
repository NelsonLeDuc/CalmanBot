package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/gif"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type RootSearch struct {
	Kind    string       `json:"kind"`
	Items   []SearchItem `json:"items"`
	Queries Queries      `json:"queries"`
}

type Queries struct {
	NextPage []Next `json:"nextPage"`
}

type Next struct {
	StartIndex int `json:"startIndex"`
}

type SearchItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Mime  string `json:"mime"`
}

func HandleGoogleImage(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query().Get("q")
	if query == "" {
		return
	}

	googleURL := "https://www.googleapis.com/customsearch/v1?key=" + os.Getenv("gis_key") + "&alt=json&searchType=image&fileType=gif&num=10&q=" + url.QueryEscape(query)
	validLinks := []string{}

	url := googleURL
	counter := 0
	for len(validLinks) < 10 && counter < 3 {
		root := googleQuery(url)
		if root == nil {
			return
		}

		validFromRoot := validLinksFromRoot(*root)
		validLinks = append(validLinks, validFromRoot...)

		url = googleURL + "&start=" + strconv.Itoa(root.Queries.NextPage[0].StartIndex)
		counter++
	}

	sliceIndex := 10
	if sliceIndex > len(validLinks) {
		sliceIndex = len(validLinks)
	}

	json, _ := json.Marshal(validLinks[0:sliceIndex])
	w.Write(json)
}

func validLinksFromRoot(root RootSearch) []string {
	validLinks := []string{}
	for _, item := range root.Items {
		if isValidGIF(item.Link) {
			validLinks = append(validLinks, item.Link)
		}
	}

	return validLinks
}

func googleQuery(url string) *RootSearch {
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var root RootSearch
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &root)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &root
}

func isValidGIF(url string) bool {

	if !isValidHTTPURLString(url) {
		return false
	}

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	byteReader := bytes.NewReader(bodyBytes)

	decode, err := gif.DecodeAll(byteReader)
	if err != nil {
		return false
	}

	return len(decode.Image) > 1
}

func isValidHTTPURLString(s string) bool {
	URL, _ := url.Parse(s)
	return (URL.Scheme == "http" || URL.Scheme == "https")
}
