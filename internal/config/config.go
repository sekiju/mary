package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/yaml.v2"
	"os"
	"runtime"
)

type Func func(*Config)

func defaultConfig() Config {
	return Config{
		Settings: Settings{
			Debug:             DebugSettings{Enable: false, Url: ""},
			OutputPath:        "output/",
			ClearOutputFolder: true,
			Threads:           runtime.NumCPU(),
		},
		Sites: map[string]SiteConfig{},
	}
}

var State Config

func Load(opts ...Func) (Config, error) {
	config := defaultConfig()
	for _, fn := range opts {
		fn(&config)
	}

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return config, nil
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	err = config.validate()
	if err != nil {
		return config, err
	}

	State = config

	return config, nil
}

func (c *Config) validate() error {
	return validation.ValidateStruct(&c.Settings,
		validation.Field(&c.Settings.Threads, validation.Required, validation.Min(1), validation.Max(runtime.NumCPU())),
	)
}
