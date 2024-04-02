package config

import (
	"errors"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"runtime"
	"strings"
)

// Config is the global config.
var Config = defaultConfig()

func defaultConfig() *Configurator {
	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: true,
	})

	return &Configurator{
		k: k,
		Settings: &settings{
			Debug:             &debugSettings{Url: ""},
			OutputPath:        "output/",
			ClearOutputFolder: true,
			Threads:           runtime.NumCPU(),
		},
		Sites: map[string]*siteConfig{},
	}
}

func (c Configurator) Load() error {
	_, err := os.Stat("config.yaml")
	if !errors.Is(err, os.ErrNotExist) {
		if err = c.k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
			return err
		}
	}

	if err = c.k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil); err != nil {
		return err
	}

	return c.k.Unmarshal("", &Config)
}
