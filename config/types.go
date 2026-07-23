package config

type config struct {
	File    FileConfig
	Runtime RuntimeFlags
}

// FileConfig holds every setting unmarshaled from the config file via koanf.
type FileConfig struct {
	Application application     `koanf:"application"`
	Output      output          `koanf:"output"`
	Sites       map[string]site `koanf:"site"`
}

// RuntimeFlags holds CLI-only state populated from flags/args, never unmarshaled by koanf.
type RuntimeFlags struct {
	PrimaryCookie    *string
	ListChaptersMode bool
	DownloadChapters []string
}

type application struct {
	CheckUpdates         bool `koanf:"check_updates"`
	MaxParallelDownloads int  `koanf:"max_parallel_downloads"`
	MaxParallelChapters  int  `koanf:"max_parallel_chapters"`
}

type output struct {
	Directory    string           `koanf:"directory"`
	CleanOnStart bool             `koanf:"clean_on_start"`
	FileFormat   OutputFileFormat `koanf:"file_format"`
}

type site struct {
	Cookie *string `koanf:"cookie"`
}

type OutputFileFormat string

const (
	AutoOutputFormat OutputFileFormat = "auto"
	PngOutputFormat  OutputFileFormat = "png"
	JpegOutputFormat OutputFileFormat = "jpeg"
	AvifOutputFormat OutputFileFormat = "avif"
	WebpOutputFormat OutputFileFormat = "webp"
)
