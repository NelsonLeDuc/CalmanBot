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
	"sync"

	"github.com/nelsonleduc/calmanbot/utility"
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

		var startIndex int
		if len(root.Queries.NextPage) > 0 {
			startIndex = root.Queries.NextPage[0].StartIndex
		} else {
			break
		}

		url = googleURL + "&start=" + strconv.Itoa(startIndex)
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
	var wg sync.WaitGroup
	out := make(chan string)

	for _, item := range root.Items {
		wg.Add(1)
		go func(link string) {
			if isValidGIF(link) {
				out <- link
			} else {
				out <- ""
			}
			wg.Done()
		}(item.Link)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	validLinks := []string{}
	for result := range out {
		if result != "" {
			validLinks = append(validLinks, result)
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

	if !utility.IsValidHTTPURLString(url) {
		return false
	}

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	byteReader := bytes.NewReader(bodyBytes)

	decode, err := gif.DecodeAll(byteReader)
	if err != nil {
		return false
	}

	return len(decode.Image) > 1
}
