package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"runtime"
	"strings"
)

var k = koanf.NewWithConf(koanf.Conf{
	Delim:       ".",
	StrictMerge: true,
})

var Data = Config{
	Settings: Settings{
		Debug:             DebugSettings{Enable: false, Url: ""},
		OutputPath:        "output/",
		ClearOutputFolder: true,
		Threads:           runtime.NumCPU(),
	},
	Sites: map[string]SiteConfig{},
}

func Load() error {
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		fmt.Printf("error loading config: %v\n", err)
	}

	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil); err != nil {
		fmt.Println(err)
	}

	err := k.Unmarshal("", &Data)
	if err != nil {
		return fmt.Errorf("error loading data: %v", err)
	}

	return nil
}
