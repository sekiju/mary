package config

import "github.com/knadh/koanf/v2"

type Configurator struct {
	k        *koanf.Koanf
	Settings *settings
	Sites    map[string]*siteConfig
}

type settings struct {
	Debug             *debugSettings
	OutputPath        string
	ClearOutputFolder bool
	Threads           int
	TargetMethod      *string
}

type debugSettings struct {
	Url string
}

type siteConfig struct {
	Session string
}
