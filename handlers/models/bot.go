package models

import (
	"strings"
)

type Bot struct {
	GroupName     string `sql:"group_name"`
	GroupID       string `sql:"group_id"`
	BotNameString string `sql:"bot_name"`
	Key           string `sql:"key"`
}

func (b Bot) BotNames() []string {
	return strings.Split(b.BotNameString, "<|>")
}

func (b Bot) SanitizedBotNames() []string {
	nameStr := strings.Replace(strings.Replace(b.BotNameString, "\\)", ")", -1), "\\(", "(", -1)
	return strings.Split(nameStr, "<|>")
}
