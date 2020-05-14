package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"

	"github.com/bwmarrin/discordgo"
	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers"
	"github.com/nelsonleduc/calmanbot/handlers/models"
	"github.com/nelsonleduc/calmanbot/service/discord"
)

// Variables used for command line parameters
var (
	token          string
	discordService discord.DSService
)

type statusTuple struct {
	gameType discordgo.GameType
	status   string
}

type dbStatus struct {
	Text string `sql:"text"`
	Type int    `sql:"type"`
}

func queryDBStatus() []statusTuple {
	queryStr := fmt.Sprintf("SELECT %s FROM discord_status", sqlstruct.Columns(dbStatus{}))
	rows, err := config.DB().Query(queryStr)
	if err != nil {
		return []statusTuple{}
	}
	defer rows.Close()

	if config.Configuration().SuperVerboseMode() {
		fmt.Println("Loaded status list:")
	}
	groupedPosts := []statusTuple{}
	for rows.Next() {
		var status dbStatus
		err := sqlstruct.Scan(&status, rows)
		if err == nil {
			var statusType discordgo.GameType
			switch status.Type {
			case 0:
				statusType = discordgo.GameTypeGame
			case 1:
				statusType = discordgo.GameTypeListening
			case 2:
				statusType = discordgo.GameTypeWatching
			default:
				continue
			}
			convertedStatus := statusTuple{statusType, status.Text}
			groupedPosts = append(groupedPosts, convertedStatus)
			if config.Configuration().SuperVerboseMode() {
				fmt.Printf("   %+v\n", convertedStatus)
			}
		}
	}

	return groupedPosts
}

func init() {
	token = os.Getenv("discord_token")
	discordService = discord.DSService{}
}

func randomStatus(excluding statusTuple) statusTuple {
	choice := excluding
	statusOptions := queryDBStatus()
	if len(statusOptions) == 0 {
		return statusTuple{discordgo.GameTypeWatching, "for questions"}
	}
	for choice == excluding {
		idx := rand.Intn(len(statusOptions))
		choice = statusOptions[idx]
	}
	return choice
}

func postStatus(s *discordgo.Session, statusTuple statusTuple) {
	log.Printf("[Tick] Setting status \"%+v\"\n", statusTuple)
	s.UpdateStatusComplex(discordgo.UpdateStatusData{
		IdleSince: nil,
		Game: &discordgo.Game{
			Name: statusTuple.status,
			Type: statusTuple.gameType,
			URL:  "",
		},
		AFK:    false,
		Status: "online"})
}

func CreateWebhook() {
	log.Println("Creating discord webhook")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalln("error opening connection,", err)
		return
	}

	status := statusTuple{}
	for ; true; <-time.Tick(30 * time.Minute) {
		status = randomStatus(status)
		postStatus(dg, status)
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := discordService.MessageFromSessionAndMessage(s, m.Message)
	monitor, _ := discordService.ServiceMonitor()
	cache := cache.NewSmartCache(monitor)

	if config.Configuration().SuperVerboseMode() {
		fmt.Printf("\n[MessageCreate fired] dssession: %+v\n", s)
		fmt.Printf("[MessageCreate fired] dsmessage: %+v\n", *m)
		fmt.Printf("[MessageCreate fired]   message: %+v\n\n", message)
	}

	handlers.HandleCalman(message, discordService, cache, models.PostGresRepo())
}
