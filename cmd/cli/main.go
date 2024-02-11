package main

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/internal/registry"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
	err := run()
	if err != nil {
		log.Error().Msgf("%v", err)
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

	var arg string
	if len(os.Args) < 2 {
		if config.Data.Settings.Debug.Enable && len(config.Data.Settings.Debug.Url) > 0 {
			arg = config.Data.Settings.Debug.Url
			log.Trace().Msgf("debug url: %s\n", arg)
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

	reader, err := registry.Default.FindReaderByDomain(uri.Hostname())
	if err != nil {
		return err
	}

	log.Info().Msgf("domain: %s | threads: %d", reader.Context().Domain, config.Data.Settings.Threads)

	startTime := time.Now()
	imageChan := make(chan connectors.ReaderImage)
	wg := &sync.WaitGroup{}

	if config.Data.Settings.ClearOutputFolder {
		err := os.RemoveAll(config.Data.Settings.OutputPath)
		if err != nil {
			return err
		}
	}

	go func() {
		err = reader.Pages(*uri, imageChan)
		if err != nil {
			fmt.Println(err)
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
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	log.Info().Msgf("Elapsed time %v\n", elapsedTime)

	return nil
}

func worker(outputPath string, imageChan <-chan connectors.ReaderImage, wg *sync.WaitGroup) error {
	defer wg.Done()

	for i := range imageChan {
		file, err := os.Create(filepath.Join(outputPath, i.FileName))
		if err != nil {
			return fmt.Errorf("failed to create file: %s", err)
		}

		img, err := i.Image()

		ext := filepath.Ext(i.FileName)
		if ext == "" {
			return fmt.Errorf("file has no extension: %s", i.FileName)
		}

		switch ext {
		case ".jpg", ".jpeg":
			err := jpeg.Encode(file, img, nil)
			if err != nil {
				return fmt.Errorf("failed to jpeg encode: %s", err)
			}
		case ".png":
			err := png.Encode(file, img)
			if err != nil {
				return fmt.Errorf("failed to png encode: %s", err)
			}
		default:
			return fmt.Errorf("unsupported image format: %s", ext)
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
