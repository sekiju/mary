package main

import (
	"559/internal/readers"
	"559/internal/registry"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/spf13/cobra"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	"path/filepath"
)

func main() {
	root.Execute()
}

type Settings struct {
	OutputPath         string `json:"output_path"`
	Url                string `json:"url"`
	ClearPreviousFiles bool   `json:"clear_folder"`
}

func (s Settings) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.OutputPath),
		validation.Field(&s.Url, validation.Required),
		validation.Field(&s.ClearPreviousFiles),
	)
}

var s = Settings{}
var root = &cobra.Command{
	Use: "559",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			s.Url = args[0]
		} else {
			return fmt.Errorf("url: cannot be blank, usage: 559.exe https://shonenjumpplus.com/episode/14079602755570530203")
		}

		if err := s.Validate(); err != nil {
			return err
		}

		return run()
	},
}

func init() {
	root.Flags().StringVarP(&s.OutputPath, "output", "o", "output", "Output folder")
	root.Flags().BoolVar(&s.ClearPreviousFiles, "clear-files", true, "Clear files in output folder")
}

func run() error {
	uri, err := url.Parse(s.Url)
	if err != nil {
		return err
	}

	reader, err := registry.Default.FindParserByDomain(uri.Hostname())
	if err != nil {
		return err
	}

	fmt.Println(reader.Context().Domain)

	imageChan := make(chan readers.ReaderImage)

	go func() {
		defer close(imageChan)

		err = reader.Pages(*uri, imageChan)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	if s.ClearPreviousFiles {
		err := os.RemoveAll(s.OutputPath)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(s.OutputPath, os.ModePerm)
	if err != nil {
		return err
	}

	for page := range imageChan {
		file, err := os.Create(filepath.Join(s.OutputPath, page.FileName))
		if err != nil {
			return fmt.Errorf("failed to create file: %s", err)
		}

		img, err := page.Image()
		if err != nil {
			return fmt.Errorf("failed to get image: %s", err)
		}

		ext := filepath.Ext(page.FileName)
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
