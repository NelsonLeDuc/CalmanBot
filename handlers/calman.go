package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/utility"
	"github.com/nelsonleduc/calmanbot/service"
	_ "github.com/nelsonleduc/calmanbot/groupme"
)

var messageService service.Service

func init() {
	messageService = *service.NewService("groupme")
}

func HandleCalman(w http.ResponseWriter, r *http.Request) {

	message := messageService.MessageFromJSON(r.Body)
	bot, _ := models.FetchBot(message.GroupID())

	if len(message.Text()) < 1 || !strings.HasPrefix(strings.ToLower(message.Text()[1:]), strings.ToLower(bot.BotName)) {
		return
	}

	actions, _ := models.FetchActions(true)
	sort.Sort(models.ByPriority(actions))

	var (
		act    models.Action
		sMatch string
	)
	for _, a := range actions {
		r, _ := regexp.Compile("(?i)" + *a.Pattern)
		matched := r.FindStringSubmatch(message.Text())
		if len(matched) > 1 && matched[1] != "" {
			sMatch = matched[1]
			act = a
			break
		}
	}

	postString := ""
	for {
		updateAction(&act, sMatch)
		if act.IsURLType() {
			postString = handleURLAction(act, w, bot)
		} else {
			postString = act.Content
		}

		if postString != "" || act.FallbackAction == nil {
			break
		} else {
			act, _ = models.FetchAction(*act.FallbackAction)
		}
	}

	if postString != "" {
		fmt.Printf("Action: %v\n", act.Content)
		fmt.Printf("Posting: %v\n", postString)
		messageService.PostText(bot.Key, postString)
	}
}

func handleURLAction(a models.Action, w http.ResponseWriter, b models.Bot) string {

	fmt.Fprintln(w, a)
	resp, err := http.Get(a.Content)

	if err == nil {

		content, _ := ioutil.ReadAll(resp.Body)
		pathString := *a.DataPath

		str := utility.ParseJSON(content, pathString)
		if str == "" {
			return ""
		} else {

			if !utility.ValidateURL(str, a.IsImageType()) {
				fmt.Printf("Invalid URL: %v\n", str)

				oldStr := str
				for i := 0; i < 3 && oldStr == str; i++ {
					str = utility.ParseJSON(content, pathString)
				}

				if !utility.ValidateURL(str, a.IsImageType()) {
					return ""
				} else {
					return str
				}
			} else {
				return str
			}
		}
	} else {
		return ""
	}

	resp.Body.Close()
	return ""
}

func updateAction(a *models.Action, text string) {
	text = url.QueryEscape(text)

	a.Content = strings.Replace(a.Content, "{_text_}", text, -1)
}
