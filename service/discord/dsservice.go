package discord

import (
	"errors"

	"github.com/bwmarrin/discordgo"

	"github.com/nelsonleduc/calmanbot/service"
)

type DSService struct{}

func (d DSService) Post(post service.Post, groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)
	if post.Type == service.PostTypeText {
		discordMessage.session.ChannelMessageSend(discordMessage.ChannelID, post.Text)
	} else if post.Type == service.PostTypeImage {
		discordMessage.session.ChannelMessageSendEmbed(discordMessage.ChannelID, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: post.Text,
			},
		})
	}
}

func (d DSService) ServiceMonitor() (service.Monitor, error) {
	return nil, errors.New("Unsupported")
}

func (d DSService) MessageFromSessionAndMessage(session *discordgo.Session, message *discordgo.Message) service.Message {
	return dsMessage{message, session}
}
