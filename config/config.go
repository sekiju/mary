package config

import (
	"errors"
	"os"

	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

var Params = config{
	File: FileConfig{
		Application: application{
			CheckUpdates:         true,
			MaxParallelDownloads: 4,
		},
		Output: output{
			Directory:    "downloads",
			CleanOnStart: false,
			FileFormat:   AutoOutputFormat,
		},
	},
}

func Load(filepath string) {
	k := koanf.NewWithConf(koanf.Conf{
		Delim:       ".",
		StrictMerge: false,
	})

	if err := k.Load(file.Provider(filepath), hcl.Parser(true)); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal().Err(err).Send()
		} else {
			log.Info().Str("filepath", filepath).Msg("Config doesn't exist. Using default configuration")
		}
	}

	if err := k.Unmarshal("", &Params.File); err != nil {
		log.Fatal().Err(err).Send()
	}
}
