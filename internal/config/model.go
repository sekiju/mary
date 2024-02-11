package config

type Config struct {
	Settings Settings
	Sites    map[string]SiteConfig
}

type Settings struct {
	Debug             DebugSettings
	OutputPath        string
	ClearOutputFolder bool
	Threads           int
}

type DebugSettings struct {
	Enable bool
	Url    string
}

type SiteConfig struct {
	Session           string
	PurchaseFreeBooks bool
}
