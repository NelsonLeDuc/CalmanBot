package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type dsMessage struct {
	*discordgo.Message
	session *discordgo.Session
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
	modifiedText := d.Content
	for _, mention := range d.Mentions {
		replaceStr := "<@" + mention.ID + ">"
		modifiedText = strings.Replace(modifiedText, replaceStr, "@"+mention.Username, -1)
	}
	return modifiedText
}

func (d dsMessage) UserType() string {
	if d.Type == discordgo.MessageTypeDefault {
		return "user"
	} else {
		return "other"
	}
}
