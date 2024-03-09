package config

import "github.com/knadh/koanf/v2"

type Configurator struct {
	k        *koanf.Koanf
	Settings Settings
	Sites    map[string]SiteConfig
}

type Settings struct {
	Debug             DebugSettings
	OutputPath        string
	ClearOutputFolder bool
	Threads           int
	TargetMethod      *string
}

type DebugSettings struct {
	Enable bool
	Url    string
}

type SiteConfig struct {
	Session           string
	PurchaseFreeBooks bool
}
