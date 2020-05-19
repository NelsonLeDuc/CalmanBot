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

	"github.com/nelsonleduc/calmanbot/config"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
	"github.com/nelsonleduc/calmanbot/utility"
)

func HandleCalman(message service.Message, providedService service.Service, cache cache.QueryCache, repo models.Repo) {
	verboseLog := config.Configuration().VerboseMode()

	if verboseLog {
		fmt.Printf("[BEGIN] Message came in: \"%v\"\n", message.Text())
	}

	if message.UserType() != "user" {
		if verboseLog {
			fmt.Println("Not from user -- aborting!")
		}
		return
	}

	bot, _ := repo.FetchBot(message.GroupID())

	if verboseLog {
		fmt.Printf("Fetched bot %#v\n", bot)
		fmt.Printf("---> Names %#v\n", bot.SanitizedBotNames())
	}

	// Make sure the message has the bot's name with a preceeding character, and that it isn't escaped
	index := -1
	for _, name := range bot.SanitizedBotNames() {
		nameIndex := strings.Index(strings.ToLower(message.Text()), strings.ToLower(name))
		if nameIndex != -1 {
			index = nameIndex
		}
	}
	isEscaped := (index >= 2 && message.Text()[index-2] == '\\')
	if isEscaped {
		if verboseLog {
			fmt.Println("Usable name not found -- aborting!")
		}
		return
	}

	var (
		postString    string
		rawPostString string
		act           models.Action
	)

	if handled, result := processBuiltins(providedService, message, bot, cache, repo); handled {
		providedService.NoteProcessing(message)
		postString = result
	} else {
		if cached := cache.CachedResponse(message.Text()); cached != nil {
			postString = *cached
		} else {
			postString, act = responseForMessage(providedService, message, bot, repo)
		}

		rawPostString = postString
		postString = updatedPostText(act, postString)
		postString = utility.ProcessedString(postString)
	}

	cacheID := cache.CacheQuery(message.Text(), postString)

	fmt.Printf("Query: %v\n", message.Text())
	if postString != "" {
		postType := postTypeForAction(act)
		fmt.Printf("Action: %v\n", act.Content)
		fmt.Printf("Type: %v\n", postType)
		fmt.Printf("Posting: %v\n\n", postString)
		providedService.Post(service.Post{bot.Key, postString, rawPostString, postType, cacheID}, message)
	} else if verboseLog {
		fmt.Print("Empty post string -- aborting!\n\n")
	}
}

func processBuiltins(service service.Service, message service.Message, bot models.Bot, cache cache.QueryCache, repo models.Repo) (bool, string) {
	verboseMode := config.Configuration().VerboseMode()
	if verboseMode {
		fmt.Println("Running builtins")
	}
	for _, b := range builtins {
		if verboseMode {
			fmt.Printf("   Check \"%v\"\n", b.trigger)
		}
		for _, name := range bot.BotNames() {
			reg, _ := regexp.Compile("(?i)@" + name + " * !" + b.trigger)
			matched := reg.FindStringSubmatch(message.Text())
			if len(matched) > 1 && matched[1] != "" {
				if verboseMode {
					fmt.Printf("      Matched \"%+v\"\n", reg)
					fmt.Printf("      Name \"%v\"\n", name)
				}
				return true, b.handler(matched, builtInParams{bot, cache, repo, service})
			}
		}
	}

	if verboseMode {
		fmt.Println("   None match")
	}
	return false, ""
}

func responseForMessage(service service.Service, message service.Message, bot models.Bot, repo models.Repo) (string, models.Action) {
	triggerHandler, triggerErr := service.ServiceTriggerWrangler()
	actions, _ := repo.FetchActions(true, triggerErr == nil)
	sort.Sort(models.ByPriority(actions))

	verboseMode := config.Configuration().VerboseMode()
	if verboseMode {
		fmt.Println("Running actions")
	}

	var (
		act    models.Action
		sMatch string
	)
	found := false
	for _, a := range actions {
		if verboseMode {
			fmt.Printf("   Check \"%v\"\n", *a.Pattern)
		}
		for _, name := range bot.BotNames() {
			regexString := strings.Replace(*a.Pattern, "{_botname_}", name+" *", -1)
			r, _ := regexp.Compile("(?i)" + regexString)
			matched := r.FindStringSubmatch(message.Text())
			if len(matched) > 1 && matched[1] != "" {
				if verboseMode {
					fmt.Printf("      Matched \"%+v\"\n", r)
					fmt.Printf("      Name \"%v\"\n", name)
				}
				sMatch = matched[1]
				act = a
				found = true
				goto Exit
			}
		}
	}
Exit:

	if found && act.WantsImmediateNote() {
		service.NoteProcessing(message)
	}

	var (
		postString string
		err        error
	)
	for {
		if verboseMode {
			fmt.Printf("                 Start: \"%v\"\n", act.Content)
		}
		updateAction(&act, sMatch)
		if verboseMode {
			fmt.Printf("                Update: \"%v\"\n", act.Content)
		}
		if act.IsURLType() {
			postString, err = handleURLAction(act, triggerHandler, bot, message)
		} else if act.IsTriggerType() {
			postString, err = handleTriggerAction(act, triggerHandler, message)
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

	if len(postString) > 0 && act.WantsPostingNote() {
		service.NoteProcessing(message)
	}

	return postString, act
}

func handleTriggerAction(action models.Action, triggerHandler service.TriggerWrangler, message service.Message) (string, error) {
	enableAction := strings.HasSuffix(action.ContentType, "ENABLE")
	disableAction := strings.HasSuffix(action.ContentType, "DISABLE")
	statusAction := strings.HasSuffix(action.ContentType, "STATUS")
	if !enableAction && !disableAction && !statusAction {
		return "", errors.New("Invalid trigger type")
	}

	triggerName := action.Content
	r, _ := regexp.Compile("(?i){_trigger\\(([^\\)]+)\\)_}")
	matched := r.FindStringSubmatch(action.Content)
	if len(matched) >= 2 {
		triggerName = matched[1]
	}

	if enableAction {
		triggerHandler.EnableTrigger(triggerName, message)
		return "Enabled", nil
	} else if statusAction {
		if triggerHandler.IsTriggerConfigured(triggerName, message) {
			return "Enabled", nil
		}
		return "Disabled", nil
	}

	triggerHandler.DisableTrigger(triggerName, message)
	return "Disabled", nil
}

func handleURLAction(a models.Action, triggerHandler service.TriggerWrangler, b models.Bot, message service.Message) (string, error) {
	url := a.Content
	r, _ := regexp.Compile("(?i){_url\\(([^\\)]+)\\)_}")
	matched := r.FindStringSubmatch(a.Content)
	if len(matched) >= 2 {
		url = matched[1]
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "CalmanBot/2.5.3")

	verboseMode := config.Configuration().VerboseMode()
	if verboseMode {
		fmt.Println("Sending URL request")
		fmt.Printf("   URL: %+v\n", *req)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Bad Response: %d length: %d", resp.StatusCode, len(content))
		return "", errors.New("Bad response from server")
	}

	if verboseMode {
		fmt.Printf("   Content: %d bytes\n", len(content))
	}

	if a.IsTriggerType() && triggerHandler != nil {
		return handleTriggerAction(a, triggerHandler, message)
	}

	pathString := *a.DataPath

	if verboseMode {
		fmt.Printf("   Path: %v\n", pathString)
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
	a.Content = strings.Replace(a.Content, "{_me_}", "localhost:"+config.Configuration().Port(), -1)

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

	if len(text) == 0 && a.IsURLPostType() {
		return text
	}

	var updated string
	if strings.Contains(*a.PostText, "{_text_}") {
		updated = strings.Replace(*a.PostText, "{_text_}", text, -1)
	} else {
		updated = *a.PostText
	}

	return updated
}

func postTypeForAction(a models.Action) service.PostType {
	if a.IsImageType() {
		return service.PostTypeImage
	} else if a.IsURLPostType() {
		return service.PostTypeURL
	} else {
		return service.PostTypeText
	}
}
