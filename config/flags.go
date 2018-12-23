package config

import "flag"

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
		discord := flag.Bool("d", false, "discord")
		flag.Parse()
		config = &configHolder{*verboseMode, *discord}
	}

	return config
}
