package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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

func HandleCalman(message service.Message, service service.Service, cache cache.QueryCache, repo models.Repo) {

	if message.UserType() != "user" {
		return
	}

	bot, _ := repo.FetchBot(message.GroupID())

	// Make sure the message has the bot's name with a preceeding character, and that it isn't escaped
	index := -1
	for _, name := range bot.SanitizedBotNames() {
		nameIndex := strings.Index(strings.ToLower(message.Text()), strings.ToLower(name))
		if nameIndex != -1 {
			index = nameIndex
		}
	}
	isEscaped := (index >= 2 && message.Text()[index-2] == '\\')
	if len(message.Text()) < 1 || index < 1 || isEscaped {
		return
	}

	var (
		postString string
		act        models.Action
	)

	if handled, result := processBuiltins(message, bot, cache, repo); handled {
		postString = result
	} else {
		if cached := cache.CachedResponse(message.Text()); cached != nil {
			postString = *cached
		} else {
			postString, act = responseForMessage(message, bot, repo)
		}

		postString = updatedPostText(act, postString)
		postString = utility.ProcessedString(postString)
	}

	cacheID := cache.CacheQuery(message.Text(), postString)

	fmt.Printf("Query: %v\n", message.Text())
	if postString != "" {
		fmt.Printf("Action: %v\n", act.Content)
		fmt.Printf("Posting: %v\n", postString)
		service.PostText(bot.Key, postString, cacheID, message)
	}
}

func processBuiltins(message service.Message, bot models.Bot, cache cache.QueryCache, repo models.Repo) (bool, string) {
	for _, b := range builtins {
		for _, name := range bot.BotNames() {
			reg, _ := regexp.Compile("(?i)&" + name + " * " + b.trigger)
			matched := reg.FindStringSubmatch(message.Text())
			if len(matched) > 1 && matched[1] != "" {
				return true, b.handler(matched, bot, cache, repo)
			}
		}
	}

	return false, ""
}

func responseForMessage(message service.Message, bot models.Bot, repo models.Repo) (string, models.Action) {
	actions, _ := repo.FetchActions(true)
	sort.Sort(models.ByPriority(actions))

	var (
		act    models.Action
		sMatch string
	)
	for _, a := range actions {
		for _, name := range bot.BotNames() {
			regexString := strings.Replace(*a.Pattern, "{_botname_}", name+" *", -1)
			r, _ := regexp.Compile("(?i)" + regexString)
			matched := r.FindStringSubmatch(message.Text())
			if len(matched) > 1 && matched[1] != "" {
				sMatch = matched[1]
				act = a
				goto Exit
			}
		}
	}
Exit:

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
			log.Printf("Failed: %s", act.Content)
			act, _ = repo.FetchAction(*act.FallbackAction)
			log.Printf("Fallback to: %s because %v", act.Content, err)
		}
	}

	return postString, act
}

func handleURLAction(a models.Action, b models.Bot) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", a.Content, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "CalmanBot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	pathString := *a.DataPath

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Bad Response: %d length: %d", resp.StatusCode, len(content))
	}

	str := utility.ParseJSON(content, pathString, utility.LinearProvider)
	for i := 0; i < 3; i++ {
		if !utility.ValidateURL(str, a.IsImageType()) {
			fmt.Printf("Invalid URL: %v\n", str)
			str = utility.ParseJSON(content, pathString, utility.LinearProvider)
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
