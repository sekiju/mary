package main

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/internal/static"
	"559/internal/utils"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	err := run()
	if err != nil {
		log.Error().Msgf("%v", err)
		fmt.Scanln()
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := config.Load()
	if err != nil {
		return err
	}

	if !config.Data.Settings.Debug.Enable {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	var arg string
	if len(os.Args) < 2 {
		if config.Data.Settings.Debug.Enable && len(config.Data.Settings.Debug.Url) > 0 {
			arg = config.Data.Settings.Debug.Url
			log.Trace().Msgf("debug url: %s", arg)
		} else {
			log.Info().Msg("provide the URL of the chapter viewer:")

			_, err := fmt.Scanln(&arg)
			if err != nil {
				return err
			}
		}
	} else {
		arg = os.Args[1]
	}

	uri, err := url.Parse(arg)
	if err != nil {
		return err
	}

	return parse(*uri)
}

func parse(uri url.URL) error {
	c, err := connectors.FindByDomain(uri.Hostname())
	if err != nil {
		return err
	}

	log.Info().Msgf("domain: %s | speed: %d image/s", c.Data().Domain, config.Data.Settings.Threads)

	urlType, err := c.ResolveType(uri)
	if err != nil {
		return err
	}

	log.Info().Msgf("url type: %s", urlType)

	imageChan := make(chan static.Image)
	wg := &sync.WaitGroup{}

	if config.Data.Settings.ClearOutputFolder {
		err = os.RemoveAll(config.Data.Settings.OutputPath)
		if err != nil {
			return err
		}
	}

	if urlType == "SHARED" {
		if config.Data.Settings.TargetMethod != nil {
			switch *config.Data.Settings.TargetMethod {
			case "book":
				urlType = "BOOK"
			case "chapter":
				urlType = "CHAPTER"
			default:
				urlType = "CHAPTER"
				log.Warn().Msgf("unknown target_method, using CHAPTER downloader")
			}
		} else {
			urlType = "CHAPTER"
			log.Warn().Msgf("unknown settings.target_method in config.yaml, using CHAPTER downloader")
		}
	}

	switch urlType {
	case "BOOK":
		if !c.Data().ChapterListAvailable {
			return fmt.Errorf("site don't supporting massive parsing")
		}

		return fmt.Errorf("book downloading not yet implemented")
	case "CHAPTER":
		chapter, err := c.Chapter(uri)
		if err != nil {
			return err
		}

		if chapter.Error != nil {
			return chapter.Error
		}

		go func() {
			err = c.Pages(chapter.ID, imageChan)
			if err != nil {
				log.Error().Msgf("%v", err)
				close(imageChan)
			}
		}()

		err = os.MkdirAll(config.Data.Settings.OutputPath, os.ModePerm)
		if err != nil {
			return err
		}

		for i := 0; i < config.Data.Settings.Threads; i++ {
			wg.Add(1)
			go func() {
				err := worker(config.Data.Settings.OutputPath, imageChan, wg)
				if err != nil {
					log.Error().Msgf("%v", err)
					os.Exit(1)
				}
			}()
		}
	}

	wg.Wait()
	return nil
}

func worker(outputPath string, imageChan <-chan static.Image, wg *sync.WaitGroup) error {
	defer wg.Done()

	for i := range imageChan {
		b, err := i.ImageFn()
		if err != nil {
			return err
		}

		err = utils.WriteImageBytes(filepath.Join(outputPath, i.FileName), b)
		if err != nil {
			return err
		}
	}

	return nil
}
