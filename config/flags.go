package config

import (
	"flag"
	"fmt"
	"os"
)

type ProcessConfig interface {
	VerboseMode() bool
	SuperVerboseMode() bool
	EnableDiscord() bool
}

type configHolder struct {
	verbose      bool
	superVerbose bool
	discord      bool
}

func (c configHolder) VerboseMode() bool {
	return c.verbose || c.superVerbose
}

func (c configHolder) SuperVerboseMode() bool {
	return c.superVerbose
}

func (c configHolder) EnableDiscord() bool {
	return c.discord
}

var config *configHolder

func Configuration() ProcessConfig {
	if config == nil {
		verboseModeFlag := flag.Bool("v", false, "more logging")
		superVerboseModeFlag := flag.Bool("vv", false, "more more logging")
		flag.Parse()
		discord := len(os.Getenv("discord_token")) > 0
		logLevel := os.Getenv("log_level")
		superVerboseMode := *superVerboseModeFlag || logLevel == "debug"
		verboseMode := *verboseModeFlag || superVerboseMode || logLevel == "info"
		config = &configHolder{verboseMode, superVerboseMode, discord}

		if superVerboseMode {
			fmt.Print("!!!! SUPER Verbose Logging enabled !!!!\n\n")
		} else if verboseMode {
			fmt.Print("!!!! Verbose Logging enabled !!!!\n\n")
		}
	}

	return config
}
