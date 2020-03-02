package main

import (
	"log"
	"math/rand"
	"os"
	"time"

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
	statusOptions  []statusTuple
)

type statusTuple struct {
	status   string
	gameType discordgo.GameType
}

func init() {
	token = os.Getenv("discord_token")
	discordService = discord.DSService{}
	statusOptions = []statusTuple{
		statusTuple{"to some sick beats", discordgo.GameTypeListening},
		statusTuple{"shit that doesn't suck yo", discordgo.GameTypeGame},
		statusTuple{"something I guess?", discordgo.GameTypeWatching},
		statusTuple{"whatever Zach thinks is good", discordgo.GameTypeGame},
		statusTuple{"to Zach ask \"Could this game get any worse?\"", discordgo.GameTypeListening},
		statusTuple{"nothing, I have no games :(", discordgo.GameTypeGame},
		statusTuple{"these l33t skillz", discordgo.GameTypeWatching},
		statusTuple{"to Alexa play Despacito", discordgo.GameTypeListening},
		statusTuple{"Jeff Goldblum movies for quotes", discordgo.GameTypeWatching},
		statusTuple{"jet fuel not melt steel beams", discordgo.GameTypeWatching},
		statusTuple{"jet beams not melt steel fuel", discordgo.GameTypeWatching},
		statusTuple{"jet steel not melt beams fuel", discordgo.GameTypeWatching},
		statusTuple{"steel beams melt jet fuel", discordgo.GameTypeWatching},
		statusTuple{"Gazorpazorpfield", discordgo.GameTypeWatching},
		statusTuple{"Star Citizen", discordgo.GameTypeGame},
		statusTuple{"The Tempest 2: Here we blow again", discordgo.GameTypeWatching},
	}
}

func randomStatus(excluding statusTuple) statusTuple {
	choice := excluding
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

	go func() {
		status := randomStatus(statusTuple{})
		postStatus(dg, status)
		c := time.Tick(30 * time.Minute)
		for range c {
			status = randomStatus(status)
			postStatus(dg, status)
		}
	}()
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

	handlers.HandleCalman(message, discordService, cache, models.PostGresRepo())
}
