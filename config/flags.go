package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type ProcessConfig interface {
	VerboseMode() bool
	SuperVerboseMode() bool
	EnableDiscord() bool
	EnableMinecraft() bool
	EnableSpotify() bool
	LocalSpotifyAuth() bool
	Port() string
	MonitorIntervalSeconds() int
	OverrideBotName() string
}

type configHolder struct {
	verbose                bool
	superVerbose           bool
	discord                bool
	minecraft              bool
	spotify                bool
	localSpotify           bool
	port                   string
	monitorIntervalSeconds int
	overrideBotName        string
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

func (c configHolder) EnableSpotify() bool {
	return c.spotify
}

func (c configHolder) LocalSpotifyAuth() bool {
	return c.localSpotify
}

func (c configHolder) Port() string {
	return c.port
}

func (c configHolder) MonitorIntervalSeconds() int {
	return c.monitorIntervalSeconds
}

func (c configHolder) OverrideBotName() string {
	return c.overrideBotName
}

var config *configHolder

func Configuration() ProcessConfig {
	if config == nil {
		localSpotifyMode := flag.Bool("local-spotify", false, "local spotify")
		verboseModeFlag := flag.Bool("v", false, "more logging")
		superVerboseModeFlag := flag.Bool("vv", false, "more more logging")
		flag.Parse()
		discord := len(os.Getenv("discord_token")) > 0
		spotify := len(os.Getenv("spotify_oauth")) > 0
		logLevel := os.Getenv("log_level")
		enableMinecraft := len(os.Getenv("enable_minecraft")) > 0

		superVerboseMode := *superVerboseModeFlag || logLevel == "debug"
		verboseMode := *verboseModeFlag || superVerboseMode || logLevel == "info"

		port := os.Getenv("PORT")
		if port == "" {
			port = "4000"
		}

		monitorIntervalEnv := os.Getenv("monitor_interval")
		monitorInterval := 0
		if len(monitorIntervalEnv) > 0 {
			if toInt, err := strconv.Atoi(monitorIntervalEnv); err == nil {
				monitorInterval = toInt
			}
		}

		overrideBotName := os.Getenv("bot_name")

		config = &configHolder{
			verbose:                verboseMode,
			superVerbose:           superVerboseMode,
			discord:                discord,
			minecraft:              enableMinecraft,
			spotify:                spotify,
			localSpotify:           *localSpotifyMode,
			port:                   port,
			monitorIntervalSeconds: monitorInterval,
			overrideBotName:        overrideBotName,
		}

		if superVerboseMode {
			fmt.Print("!!!! SUPER Verbose Logging enabled !!!!\n\n")
		} else if verboseMode {
			fmt.Print("!!!! Verbose Logging enabled !!!!\n\n")
		}
	}

	return config
}
