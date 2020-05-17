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
	EnableMinecraft() bool
	Port() string
}

type configHolder struct {
	verbose      bool
	superVerbose bool
	discord      bool
	minecraft    bool
	port         string
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

func (c configHolder) EnableMinecraft() bool {
	return c.minecraft
}

func (c configHolder) Port() string {
	return c.port
}

var config *configHolder

func Configuration() ProcessConfig {
	if config == nil {
		verboseModeFlag := flag.Bool("v", false, "more logging")
		superVerboseModeFlag := flag.Bool("vv", false, "more more logging")
		flag.Parse()
		discord := len(os.Getenv("discord_token")) > 0
		logLevel := os.Getenv("log_level")
		enableMinecraft := len(os.Getenv("enable_minecraft")) > 0

		superVerboseMode := *superVerboseModeFlag || logLevel == "debug"
		verboseMode := *verboseModeFlag || superVerboseMode || logLevel == "info"

		port := os.Getenv("PORT")
		if port == "" {
			port = "4000"
		}

		config = &configHolder{verboseMode, superVerboseMode, discord, enableMinecraft, port}

		if superVerboseMode {
			fmt.Print("!!!! SUPER Verbose Logging enabled !!!!\n\n")
		} else if verboseMode {
			fmt.Print("!!!! Verbose Logging enabled !!!!\n\n")
		}
	}

	return config
}
