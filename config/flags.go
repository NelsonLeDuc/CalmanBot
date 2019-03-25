package config

import (
	"flag"
	"fmt"
	"os"
)

type ProcessConfig interface {
	VerboseMode() bool
	EnableDiscord() bool
}

type configHolder struct {
	verbose bool
	discord bool
}

func (c configHolder) VerboseMode() bool {
	return c.verbose
}

func (c configHolder) EnableDiscord() bool {
	return c.discord
}

var config *configHolder

func Configuration() ProcessConfig {
	if config == nil {
		verboseMode := flag.Bool("v", false, "more logging")
		flag.Parse()
		discord := len(os.Getenv("discord_token")) > 0
		config = &configHolder{*verboseMode, discord}

		if *verboseMode {
			fmt.Print("!!!! Verbose Logging enabled !!!!\n\n")
		}
	}

	return config
}
