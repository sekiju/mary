package config

type Config struct {
	Settings Settings              `yaml:"settings"`
	Sites    map[string]SiteConfig `yaml:"sites"`
}

type Settings struct {
	EnableDebug       bool   `yaml:"enable_debug"`
	OutputPath        string `yaml:"output_path"`
	ClearOutputFolder bool   `yaml:"clear_output_folder"`
	Threads           int    `yaml:"threads"`
}

type SiteConfig struct {
	Session           string `yaml:"session"`
	PurchaseFreeBooks bool   `yaml:"purchase_free_books,omitempty"`
}
