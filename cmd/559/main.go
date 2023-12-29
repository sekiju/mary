package main

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/internal/registry"
	"fmt"
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
		fmt.Print(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: 559 <url to chapter viewer>")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	uri, err := url.Parse(os.Args[1])
	if err != nil {
		return err
	}

	reader, err := registry.Default.FindReaderByDomain(uri.Hostname())
	if err != nil {
		return err
	}

	fmt.Println(reader.Context().Domain, " | ", fmt.Sprintf("threads: %d", cfg.Settings.Threads))

	startTime := time.Now()
	imageChan := make(chan connectors.ReaderImage)
	wg := &sync.WaitGroup{}

	if cfg.Settings.ClearOutputFolder {
		err := os.RemoveAll(cfg.Settings.OutputPath)
		if err != nil {
			return err
		}
	}

	go func() {
		err = reader.Pages(*uri, imageChan)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	err = os.MkdirAll(cfg.Settings.OutputPath, os.ModePerm)
	if err != nil {
		return err
	}

	for i := 0; i < int(cfg.Settings.Threads); i++ {
		wg.Add(1)
		go func() {
			err := worker(cfg.Settings.OutputPath, imageChan, wg)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	fmt.Printf("Elapsed time %v+\n", elapsedTime)

	return nil
}

func worker(outputPath string, imageChan <-chan connectors.ReaderImage, wg *sync.WaitGroup) error {
	for i := range imageChan {
		file, err := os.Create(filepath.Join(outputPath, i.FileName))
		if err != nil {
			return fmt.Errorf("failed to create file: %s", err)
		}

		img, err := i.Image()

		ext := filepath.Ext(i.FileName)
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

	wg.Done()

	return nil
}
