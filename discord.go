package main

import (
	"fmt"
	"os"

	"github.com/nelsonleduc/calmanbot/cache"
	"github.com/nelsonleduc/calmanbot/handlers"
	"github.com/nelsonleduc/calmanbot/service/discord"

	"github.com/nelsonleduc/calmanbot/service"

	"github.com/bwmarrin/discordgo"
	"github.com/nelsonleduc/calmanbot/handlers/models"
)

// Variables used for command line parameters
var (
	token          string
	discordService discord.DSService
)

func init() {
	token = os.Getenv("discord_token")
	discordService = (*service.NewService("discord")).(discord.DSService)
}

func CreateWebhook() {
	fmt.Println("creating hook")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
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

	handlers.HandleCalman(message, discordService, cache, models.PostGresRepo())
}
