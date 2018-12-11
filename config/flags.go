package config

import "flag"

type ProcessConfig interface {
	VerboseMode() bool
}

type configHolder struct {
	verbose bool
}

func (c configHolder) VerboseMode() bool {
	return c.verbose
}

var config *configHolder

func Configuration() ProcessConfig {
	if config == nil {
		verboseMode := flag.Bool("v", false, "more logging")
		flag.Parse()
		config = &configHolder{*verboseMode}
	}

	return config
}
