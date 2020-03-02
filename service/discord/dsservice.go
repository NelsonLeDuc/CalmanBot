package discord

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/bwmarrin/discordgo"

	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
)

type DSService struct{}

func (d DSService) Post(post service.Post, groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)
	if post.Type == service.PostTypeText {
		discordMessage.session.ChannelMessageSend(discordMessage.ChannelID, post.Text)
	} else if post.Type == service.PostTypeImage {
		var footer *discordgo.MessageEmbedFooter
		if postURL, err := url.Parse(post.Text); err == nil {
			footer = &discordgo.MessageEmbedFooter{
				Text: postURL.Host,
			}
		} else {
			footer = nil
		}
		discordMessage.session.ChannelMessageSendEmbed(discordMessage.ChannelID, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: post.Text,
			},
			Footer: footer,
		})
	}
}

func (d DSService) ServiceMonitor() (service.Monitor, error) {
	return nil, errors.New("Unsupported")
}

func (d DSService) MessageFromSessionAndMessage(session *discordgo.Session, message *discordgo.Message) service.Message {
	processed := processText(session, message)
	return dsMessage{message, session, processed}
}

func processText(session *discordgo.Session, message *discordgo.Message) string {
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

	if verboseLog {
		fmt.Printf("Final discord text: \"%s\"\n", modifiedText)
	}

	return modifiedText
}
