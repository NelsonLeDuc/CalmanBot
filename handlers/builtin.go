package handlers

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
)

type builtinDescription struct {
	trigger     string
	description string
}

type builtin struct {
	builtinDescription
	handler func([]string, models.Bot, cache.QueryCache) string
}

// Descriptions
var helpDescription = builtinDescription{
	"(help)",
	"See this",
}
var topDescription = builtinDescription{
	"(top)",
	"List top 5 liked posts",
}
var showDescription = builtinDescription{
	"show ([1-5])",
	"Repost nth top post",
}

var descriptions = []builtinDescription{
	helpDescription,
	topDescription,
	showDescription,
}
var builtins = []builtin{
	builtin{
		helpDescription,
		responseForHelp,
	},
	builtin{
		topDescription,
		responseForLeaderboard,
	},
	builtin{
		showDescription,
		responseForShow,
	},
}

// Handlers
func responseForLeaderboard(matched []string, bot models.Bot, cache cache.QueryCache) string {
	entries := cache.LeaderboardEntries(bot.GroupID, 5)
	leaderboardAccumulatr := "Top posts:"
	for _, e := range entries {
		leaderboardAccumulatr += "\n" + strconv.Itoa(e.LikeCount) + "    " + e.Query
	}

	return leaderboardAccumulatr
}

func responseForShow(matched []string, bot models.Bot, cache cache.QueryCache) string {
	entries := cache.LeaderboardEntries(bot.GroupID, 5)
	num, error := strconv.Atoi(matched[1])
	num--
	if len(entries) <= num || error != nil {
		return "There is nothing to display"
	}

	return entries[num].Result
}

func responseForHelp(matched []string, bot models.Bot, cache cache.QueryCache) string {
	actions, _ := models.FetchActions(true)
	sort.Sort(models.ByPriority(actions))

	helpAccumulator := "Commands:"
	longest := 0
	for _, a := range actions {
		if a.Description == nil || a.Pattern == nil {
			continue
		}

		if len(*a.Pattern) > longest {
			longest = len(*a.Pattern)
		}
	}
	for _, b := range descriptions {
		length := len("&" + bot.BotName + " " + b.trigger)
		if length > longest {
			longest = length
		}
	}
	paddingFmt := fmt.Sprintf("%%-%ds", longest+2)

	for _, a := range actions {
		if a.Description == nil || a.Pattern == nil {
			continue
		}
		printablePattern := *a.Pattern
		printablePattern = strings.Replace(printablePattern, "{_botname_}", bot.BotName, -1)
		re := regexp.MustCompile("^\\[(.)\\]")
		matched := re.FindStringSubmatch(printablePattern)
		thing := ""
		if len(matched) > 1 && matched[1] != "" {
			thing = matched[1]
		}
		printablePattern = re.ReplaceAllLiteralString(printablePattern, thing)

		helpAccumulator += "\n" + fmt.Sprintf(paddingFmt, "\""+printablePattern+"\"") + "\n\t" + *a.Description
	}
	for _, b := range descriptions {
		printablePattern := "&" + bot.BotName + " " + b.trigger
		helpAccumulator += "\n" + fmt.Sprintf(paddingFmt, "\""+printablePattern+"\"") + "\n\t" + b.description
	}

	return helpAccumulator
}
