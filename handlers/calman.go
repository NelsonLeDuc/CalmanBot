package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
	"github.com/nelsonleduc/calmanbot/utility"
)

func HandleCalman(message service.Message, service service.Service, cache cache.QueryCache) {

	if message.UserType() != "user" {
		return
	}

	bot, _ := models.FetchBot(message.GroupID())

	// Make sure the message has the bot's name with a preceeding character, and that it isn't escaped
	index := strings.Index(strings.ToLower(message.Text()), strings.ToLower(bot.BotName))
	isEscaped := (index >= 2 && message.Text()[index-2] == '\\')
	if len(message.Text()) < 1 || index < 1 || isEscaped {
		return
	}

	var (
		postString string
		act        models.Action
	)

	if cached := cache.CachedResponse(message.Text()); cached != nil {
		postString = *cached
	} else {
		postString, act = responseForMessage(message, bot)
	}

	postString = updatedPostText(act, postString)
	postString = utility.ProcessedString(postString)

	cacheID := cache.CacheQuery(message.Text(), postString)

	fmt.Printf("Query: %v\n", message.Text())
	if postString != "" {
		fmt.Printf("Action: %v\n", act.Content)
		fmt.Printf("Posting: %v\n", postString)
		service.PostText(bot.Key, postString, cacheID, message)
	}
}

func responseForMessage(message service.Message, bot models.Bot) (string, models.Action) {
	actions, _ := models.FetchActions(true)
	sort.Sort(models.ByPriority(actions))

	var (
		act    models.Action
		sMatch string
	)
	for _, a := range actions {
		regexString := strings.Replace(*a.Pattern, "{_botname_}", bot.BotName, -1)
		r, _ := regexp.Compile("(?i)" + regexString)
		matched := r.FindStringSubmatch(message.Text())
		if len(matched) > 1 && matched[1] != "" {
			sMatch = matched[1]
			act = a
			break
		}
	}

	var (
		postString string
		err        error
	)
	for {
		updateAction(&act, sMatch)
		if act.IsURLType() {
			postString, err = handleURLAction(act, bot)
		} else {
			postString = act.Content
		}

		if (err == nil && postString != "") || act.FallbackAction == nil {
			break
		} else {
			act, _ = models.FetchAction(*act.FallbackAction)
		}
	}

	return postString, act
}

func handleURLAction(a models.Action, b models.Bot) (string, error) {

	resp, err := http.Get(a.Content)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	content, _ := ioutil.ReadAll(resp.Body)
	pathString := *a.DataPath

	str := utility.ParseJSON(content, pathString)
	for i := 0; i < 3; i++ {
		if !utility.ValidateURL(str, a.IsImageType()) {
			fmt.Printf("Invalid URL: %v\n", str)
			str = utility.ParseJSON(content, pathString)
		} else {
			return str, nil
		}
	}

	return "", errors.New("Failed to handle URL action")
}

func updateAction(a *models.Action, text string) {
	text = url.QueryEscape(text)

	a.Content = strings.Replace(a.Content, "{_text_}", text, -1)

	r, _ := regexp.Compile("(?i){_key\\((.+)\\)_}")
	matched := r.FindStringSubmatch(a.Content)
	if len(matched) >= 2 {
		keyVal := os.Getenv(matched[1] + "_key")
		a.Content = strings.Replace(a.Content, matched[0], keyVal, -1)
	}
}

func updatedPostText(a models.Action, text string) string {
	if a.PostText == nil {
		return text
	}

	var updated string
	if strings.Contains(*a.PostText, "{_text_}") {
		updated = strings.Replace(*a.PostText, "{_text_}", text, -1)
	} else {
		updated = *a.PostText + text
	}

	return updated
}
