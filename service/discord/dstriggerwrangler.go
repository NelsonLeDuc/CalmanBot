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

func (t triggerWrangler) EnableTrigger(id string, groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)

	queryStr := "INSERT INTO discord_triggers(channel_id, trigger_id) VALUES($1, $2) ON CONFLICT DO NOTHING"
	config.DB().Exec(queryStr, discordMessage.ChannelID, id)
}

func (t triggerWrangler) DisableTrigger(id string, groupMessage service.Message) {
	discordMessage := groupMessage.(dsMessage)

	queryStr := "DELETE FROM discord_triggers WHERE channel_id = $1 AND trigger_id = $2"
	config.DB().Exec(queryStr, discordMessage.ChannelID, id)
}

func (t triggerWrangler) IsTriggerConfigured(id string, groupMessage service.Message) bool {
	discordMessage := groupMessage.(dsMessage)

	queryStr := "SELECT count(*) FROM discord_triggers WHERE channel_id = $1 AND trigger_id = $2"
	row := config.DB().QueryRow(queryStr, discordMessage.ChannelID, id)
	var count int
	err := row.Scan(&count)
	fmt.Printf("%v\n", err)
	if err != nil {
		return false
	}
	fmt.Printf("%v\n", count)

	return count > 0
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
		fmt.Printf("   %+v\n", trigger)
		t.service.postToChannel(post, trigger.ChannelID)
	}
}
