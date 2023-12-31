package config

type Config struct {
	Settings Settings              `yaml:"settings"`
	Sites    map[string]SiteConfig `yaml:"sites"`
}

type Settings struct {
	Debug             DebugSettings `yaml:"debug"`
	OutputPath        string        `yaml:"output_path"`
	ClearOutputFolder bool          `yaml:"clear_output_folder"`
	Threads           int           `yaml:"threads"`
}

type DebugSettings struct {
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
}

type SiteConfig struct {
	Session           string `yaml:"session"`
	PurchaseFreeBooks bool   `yaml:"purchase_free_books,omitempty"`
}
