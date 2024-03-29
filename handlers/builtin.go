package handlers

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service"
)

const currentCalmanBotVersion string = "v2.16.1"

type builtinDescription struct {
	trigger     string
	description string
}

type builtInParams struct {
	bot     models.Bot
	cache   cache.QueryCache
	repo    models.Repo
	service service.Service
}

type builtin struct {
	builtinDescription
	featureChecker func(service.Service) bool
	handler        func([]string, builtInParams) string
}

// Descriptions
var helpDescription = builtinDescription{
	"(help)",
	"See this",
}
var topDescription = builtinDescription{
	"(top)",
	"List top 10 liked posts",
}
var showDescription = builtinDescription{
	"show (10|[1-9])(?:$| )",
	"Repost nth top post",
}

var versionDescription = builtinDescription{
	"(version)",
	"Display current version",
}

var builtins []builtin

func init() {
	builtins = []builtin{
		{
			helpDescription,
			noFeatureCheck,
			responseForHelp,
		},
		{
			topDescription,
			leaderBoardFeatureCheck,
			responseForLeaderboard,
		},
		{
			showDescription,
			leaderBoardFeatureCheck,
			responseForShow,
		},
		{
			versionDescription,
			noFeatureCheck,
			responseForVersion,
		},
	}
}

// Handlers
func responseForLeaderboard(matched []string, params builtInParams) string {
	entries := params.cache.LeaderboardEntries(params.bot.GroupID, 10)
	leaderboardAccumulatr := "Top posts:"
	for _, e := range entries {
		leaderboardAccumulatr += "\n" + strconv.Itoa(e.LikeCount) + "    " + e.Query
	}

	return leaderboardAccumulatr
}

func responseForShow(matched []string, params builtInParams) string {
	entries := params.cache.LeaderboardEntries(params.bot.GroupID, 10)
	num, error := strconv.Atoi(matched[1])
	num--
	if len(entries) <= num || error != nil {
		return "There is nothing to display"
	}

	return entries[num].Result
}

func responseForVersion(matched []string, params builtInParams) string {
	return "I'm currently running " + currentCalmanBotVersion
}

func responseForHelp(matched []string, params builtInParams) string {
	_, e := params.service.ServiceTriggerWrangler()
	actions, _ := params.repo.FetchActions(true, e == nil)
	sort.Sort(models.ByPriority(actions))
	botName := params.bot.SanitizedBotNames()[0]

	helpAccumulator := "```Commands:"
	longest := 0
	for _, a := range actions {
		if a.Description == nil || a.Pattern == nil {
			continue
		}

		if len(*a.Pattern) > longest {
			longest = len(*a.Pattern)
		}
	}
	for _, b := range builtins {
		if !b.featureChecker(params.service) {
			continue
		}
		length := len("@" + botName + " !" + b.trigger)
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
		printablePattern = strings.Replace(printablePattern, "{_botname_}", botName, -1)
		re := regexp.MustCompile(`^\[(.)\]`)
		matched := re.FindStringSubmatch(printablePattern)
		thing := ""
		if len(matched) > 1 && matched[1] != "" {
			thing = matched[1]
		}
		printablePattern = re.ReplaceAllLiteralString(printablePattern, thing)

		helpAccumulator += "\n" + fmt.Sprintf(paddingFmt, printablePattern) + "\n\t" + *a.Description
	}
	for _, b := range builtins {
		if !b.featureChecker(params.service) {
			continue
		}
		printablePattern := "@" + botName + " !" + b.trigger
		helpAccumulator += "\n" + fmt.Sprintf(paddingFmt, printablePattern) + "\n\t" + b.description
	}

	return helpAccumulator + "```"
}

// Feature checker
func noFeatureCheck(s service.Service) bool {
	return true
}

func leaderBoardFeatureCheck(s service.Service) bool {
	return s.SupportsBuiltInFeature(service.BuiltInFeatureLeaderboard)
}
