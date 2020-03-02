package discord

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/nelsonleduc/calmanbot/config"
)

type dsMessage struct {
	*discordgo.Message
	session *discordgo.Session

	processedText *string
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
	if d.processedText != nil {
		return *d.processedText
	}
	verboseLog := config.Configuration().VerboseMode()
	modifiedText := d.Content
	if verboseLog {
		fmt.Printf("Parsing discord text: \"%s\"\n", modifiedText)
	}
	for _, mention := range d.Mentions {
		re := regexp.MustCompile("<@.?" + mention.ID + ">")
		modifiedText = re.ReplaceAllString(modifiedText, "@"+mention.Username)
		if verboseLog {
			fmt.Printf("Replacing \"%v\" with \"%v\"\n", mention.ID, mention.Username)
		}
	}
	if verboseLog {
		fmt.Printf("Final discord text: \"%s\"\n", modifiedText)
	}
	d.processedText = &modifiedText
	return modifiedText
}

func (d dsMessage) UserType() string {
	if d.Type == discordgo.MessageTypeDefault {
		return "user"
	} else {
		return "other"
	}
}
