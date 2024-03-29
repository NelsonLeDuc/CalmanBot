package discord

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
)

type DSService interface {
	service.Service
	MessageFromSessionAndMessage(session *discordgo.Session, message *discordgo.Message) service.Message
}
type dsService struct {
	session *discordgo.Session
}

func NewDSService(discordSession *discordgo.Session) DSService {
	newService := dsService{discordSession}
	service.RegisterServiceForTriggers(newService)
	return newService
}

func (d dsService) Post(post service.Post, groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)
	d.postToChannel(post, discordMessage.ChannelID)
}

func (d dsService) postToChannel(post service.Post, channelID string) {
	if post.Type == service.PostTypeText {
		d.session.ChannelMessageSend(channelID, post.Text)
	} else if post.Type == service.PostTypeImage {
		var footer *discordgo.MessageEmbedFooter
		if postURL, err := url.Parse(post.Text); err == nil {
			footer = &discordgo.MessageEmbedFooter{
				Text: postURL.Host,
			}
		} else {
			footer = nil
		}
		d.session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: post.Text,
			},
			Footer: footer,
		})
	} else if post.Type == service.PostTypeURL {
		d.session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
			URL:   post.RawText,
			Title: post.Text,
		})
	}
}

func (d dsService) NoteProcessing(groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)

	emojis := []string{"🍉", "🤓", "🦥", "🍁", "🌝", "🐈", "🦞", "👒", "💸", "👁️", "❤️‍🔥", "👻", "👐", "🤜", "🦑", "🌚", "💸"}
	chosen := emojis[rand.Intn(len(emojis))]
	err := discordMessage.session.MessageReactionAdd(discordMessage.ChannelID, discordMessage.ID, chosen)
	if err != nil && config.Configuration().VerboseMode() {
		fmt.Printf("Error posting reaction: chose: %v err: %v\n", chosen, err)
	}
}

func (d dsService) ServiceMonitor() (service.Monitor, error) {
	return nil, errors.New("Unsupported")
}

func (d dsService) MessageFromSessionAndMessage(session *discordgo.Session, message *discordgo.Message) service.Message {
	c, err := session.Channel(message.ChannelID)
	isDirect := false
	if err != nil && config.Configuration().VerboseMode() {
		fmt.Printf("Error fetching channel %v", err)
	} else {
		isDirect = c.Type == discordgo.ChannelTypeDM
	}

	processed := processText(session, message, isDirect)

	return dsMessage{message, session, processed}
}

func (d dsService) SupportsBuiltInFeature(feature service.BuiltInFeature) bool {
	switch feature {
	case service.BuiltInFeatureLeaderboard:
		return false
	}
	return false
}

func processText(session *discordgo.Session, message *discordgo.Message, isDirect bool) string {
	verboseLog := config.Configuration().VerboseMode()
	modifiedText := message.Content
	if verboseLog {
		fmt.Printf("Parsing discord text: \"%s\"\n", modifiedText)
	}
	for _, mention := range message.Mentions {
		re := regexp.MustCompile("<@.?" + mention.ID + ">")
		modifiedText = re.ReplaceAllString(modifiedText, "@"+mention.Username)
		if verboseLog {
			fmt.Printf("Replacing \"%v\" with \"%v\"\n", mention.ID, mention.Username)
		}
	}
	if len(message.MentionRoles) > 0 {
		g, _ := session.GuildRoles(message.GuildID)
		for _, mention := range message.MentionRoles {
			for _, role := range g {
				if role.ID == mention {
					re := regexp.MustCompile("<@.?" + mention + ">")
					modifiedText = re.ReplaceAllString(modifiedText, "@"+role.Name)
					if verboseLog {
						fmt.Printf("Replacing role \"%v\" with \"%v\"\n", mention, role.Name)
					}
				}
			}
		}
	}
	re := regexp.MustCompile("<@.?" + session.State.User.ID + ">")
	modifiedText = re.ReplaceAllString(modifiedText, "@"+session.State.User.Username)

	if isDirect && !strings.Contains(modifiedText, "@"+session.State.User.Username+" ") {
		modifiedText = "@" + session.State.User.Username + " " + modifiedText
	}

	if verboseLog {
		fmt.Printf("Final discord text: \"%s\"\n", modifiedText)
	}

	return modifiedText
}
