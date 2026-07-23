package config

const CurrentConfigVersion = 1

type config struct {
	File    FileConfig
	Runtime RuntimeFlags
}

// FileConfig holds every setting persisted to the config file.
type FileConfig struct {
	Schema      string          `json:"$schema,omitempty"`
	Version     int             `json:"version"`
	Application application     `json:"application"`
	Output      output          `json:"output"`
	Sites       map[string]Site `json:"site,omitempty"`
}

// RuntimeFlags holds CLI-only state populated from flags/args, never persisted.
type RuntimeFlags struct {
	ListChaptersMode bool
	DownloadChapters []string
}

type application struct {
	CheckUpdates         bool `json:"check_updates"`
	MaxParallelDownloads int  `json:"max_parallel_downloads"`
	MaxParallelChapters  int  `json:"max_parallel_chapters"`
}

type output struct {
	Directory    string           `json:"directory"`
	CleanOnStart bool             `json:"clean_on_start"`
	FileFormat   OutputFileFormat `json:"file_format"`
}

type Site struct {
	Cookie *string `json:"cookie,omitempty"`
}

type OutputFileFormat string

const (
	AutoOutputFormat OutputFileFormat = "auto"
	PngOutputFormat  OutputFileFormat = "png"
	JpegOutputFormat OutputFileFormat = "jpeg"
	AvifOutputFormat OutputFileFormat = "avif"
	WebpOutputFormat OutputFileFormat = "webp"
)
