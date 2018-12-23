package discord

import (
	"io"

	"github.com/bwmarrin/discordgo"

	"github.com/nelsonleduc/calmanbot/service"
)

type DSService struct{}

func init() {
	service.AddService("discord", DSService{})
}

func (d DSService) PostText(key, text string, pType service.PostType, cacheID int, groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)
	if pType == service.PostTypeText {
		discordMessage.session.ChannelMessageSend(discordMessage.ChannelID, text)
	} else if pType == service.PostTypeImage {
		discordMessage.session.ChannelMessageSendEmbed(discordMessage.ChannelID, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: text,
			},
		})
	}
}

func (d DSService) MessageFromJSON(reader io.Reader) service.Message {
	return dsMessage{}
}

func (d DSService) ServiceMonitor() (service.Monitor, error) {
	return DiscordMonitor{}, nil
}

func (d DSService) MessageFromSessionAndMessage(session *discordgo.Session, message *discordgo.Message) service.Message {
	return dsMessage{message, session}
}
