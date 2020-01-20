package discord

import (
	"errors"
	"net/url"

	"github.com/bwmarrin/discordgo"

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
	return dsMessage{message, session}
}
