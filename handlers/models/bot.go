package models

type Bot struct {
	GroupName string `sql:"group_name"`
	GroupID   string `sql:"group_id"`
	BotName   string `sql:"bot_name"`
	Key       string `sql:"key"`
}
