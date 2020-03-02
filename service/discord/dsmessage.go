package discord

import (
	"github.com/bwmarrin/discordgo"
)

type dsMessage struct {
	*discordgo.Message
	session *discordgo.Session

	processedText string
}

func (d dsMessage) GroupID() string {
	return "discord"
}

func (d dsMessage) UserName() string {
	return d.Author.Username
}

func (d dsMessage) UserID() string {
	return d.Author.ID
}

func (d dsMessage) MessageID() string {
	return d.ID
}

func (d dsMessage) Text() string {
	return d.processedText
}

func (d dsMessage) UserType() string {
	if d.Type == discordgo.MessageTypeDefault {
		return "user"
	} else {
		return "other"
	}
}
