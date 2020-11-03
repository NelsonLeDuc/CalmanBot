package discord

import (
	"fmt"

	"github.com/kisielk/sqlstruct"
	"github.com/nelsonleduc/calmanbot/config"
	"github.com/nelsonleduc/calmanbot/service"
)

type triggerWrangler struct {
	service dsService
}

type storedTrigger struct {
	TriggerID string `sql:"trigger_id"`
	ChannelID string `sql:"channel_id"`
}

func (d dsService) ServiceTriggerWrangler() (service.TriggerWrangler, error) {
	return triggerWrangler{d}, nil
}

func (t triggerWrangler) EnableTrigger(id string, groupMessage service.Message, forGuild bool) {
	discordMessage := groupMessage.(dsMessage)

	queryStr := "INSERT INTO discord_triggers(channel_id, trigger_id) VALUES($1, $2) ON CONFLICT DO NOTHING"
	if forGuild {
		config.DB().Exec(queryStr, discordMessage.GuildID, id)
	} else {
		config.DB().Exec(queryStr, discordMessage.ChannelID, id)
	}
}

func (t triggerWrangler) DisableTrigger(id string, groupMessage service.Message, forGuild bool) {
	discordMessage := groupMessage.(dsMessage)

	queryStr := "DELETE FROM discord_triggers WHERE channel_id = $1 AND trigger_id = $2"
	if forGuild {
		config.DB().Exec(queryStr, discordMessage.GuildID, id)
	} else {
		config.DB().Exec(queryStr, discordMessage.ChannelID, id)
	}
}

func (t triggerWrangler) IsTriggerConfigured(id string, groupMessage service.Message, forGuild bool) bool {
	discordMessage := groupMessage.(dsMessage)
	if forGuild {
		return t.HasTrigger(id, discordMessage.GuildID)
	}

	return t.HasTrigger(id, discordMessage.ChannelID)
}

func (t triggerWrangler) HandleTrigger(id string, post service.Post) {
	queryStr := fmt.Sprintf("SELECT %s FROM discord_triggers WHERE trigger_id = $1", sqlstruct.Columns(storedTrigger{}))
	rows, err := config.DB().Query(queryStr, id)
	if err != nil {
		return
	}
	defer rows.Close()

	if config.Configuration().SuperVerboseMode() {
		fmt.Println("Loaded trigger list:")
	}
	for rows.Next() {
		var trigger storedTrigger
		err := sqlstruct.Scan(&trigger, rows)
		if err != nil {
			continue
		}
		if config.Configuration().SuperVerboseMode() {
			fmt.Printf("   %+v\n", trigger)
		}
		t.service.postToChannel(post, trigger.ChannelID)
	}
}

func (t triggerWrangler) HasTrigger(id string, groupID string) bool {
	queryStr := "SELECT count(*) FROM discord_triggers WHERE channel_id = $1 AND trigger_id = $2"
	row := config.DB().QueryRow(queryStr, groupID, id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}
