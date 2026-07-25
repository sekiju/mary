package config

const CurrentConfigVersion = 1

type config struct {
	File    FileConfig
	Runtime RuntimeFlags
}

// FileConfig holds every setting persisted to the config file.
type FileConfig struct {
	Schema   string          `json:"$schema,omitempty"`
	Version  int             `json:"version"`
	Settings settings        `json:"settings"`
	Sites    map[string]Site `json:"site,omitempty"`
}

// RuntimeFlags holds CLI-only state populated from flags/args, never persisted.
type RuntimeFlags struct {
	ListChaptersMode bool
	DownloadChapters []string
}

type settings struct {
	CheckForUpdates             bool   `json:"checkForUpdates"`
	AutoInstallUpdates          bool   `json:"autoInstallUpdates"`
	MaxParallelPageFetches      int    `json:"maxParallelPageFetches"`
	MaxParallelChaptersDownload int    `json:"maxParallelChaptersDownload"`
	OutputDirectory             string `json:"outputDirectory"`
}

type Site struct {
	Cookie *string `json:"cookie,omitempty"`
}
