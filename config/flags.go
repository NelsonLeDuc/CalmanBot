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
	MinecraftAddress() string
}

type configHolder struct {
	verbose      bool
	superVerbose bool
	discord      bool
	minecraft    string
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

func (c configHolder) MinecraftAddress() string {
	return c.minecraft
}

var config *configHolder

func Configuration() ProcessConfig {
	if config == nil {
		verboseModeFlag := flag.Bool("v", false, "more logging")
		superVerboseModeFlag := flag.Bool("vv", false, "more more logging")
		minecraftFlag := flag.String("mc", "", "Minecraft server to monitor (including port)")
		flag.Parse()
		discord := len(os.Getenv("discord_token")) > 0
		logLevel := os.Getenv("log_level")
		minecraftEnv := os.Getenv("minecraft_address")

		superVerboseMode := *superVerboseModeFlag || logLevel == "debug"
		verboseMode := *verboseModeFlag || superVerboseMode || logLevel == "info"
		var minecraftAddress string
		if len(minecraftEnv) > 0 {
			minecraftAddress = minecraftEnv
		} else {
			minecraftAddress = *minecraftFlag
		}

		config = &configHolder{verboseMode, superVerboseMode, discord, minecraftAddress}

		if superVerboseMode {
			fmt.Print("!!!! SUPER Verbose Logging enabled !!!!\n\n")
		} else if verboseMode {
			fmt.Print("!!!! Verbose Logging enabled !!!!\n\n")
		}
	}

	return config
}
