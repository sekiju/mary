package downloader

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/knst0/mdl/config"
	"github.com/knst0/mdl/sdk/manga"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"resty.dev/v3"
)

var httpClient = resty.New()

func getReader(ctx context.Context, page *manga.Page) (io.Reader, error) {
	res, err := httpClient.R().SetContext(ctx).SetHeaders(page.Headers).Get(page.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to download page: %w", err)
	}

	b := res.Bytes()

	if page.Decode == nil {
		return bytes.NewReader(b), nil
	}

	b, err = page.Decode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to decode page: %w", err)
	}

	return bytes.NewReader(b), nil
}

func saveFile(dir, filename string, r io.Reader) error {
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil

}

func generateFilename(originalName string, format config.OutputFileFormat) string {
	return originalName[:1+len(originalName)-len(filepath.Ext(originalName))] + string(format)
}

func saveEncodedImage(dir, filename string, format config.OutputFileFormat, r io.Reader) error {
	img, _, err := image.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	filename = generateFilename(filename, format)
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	switch format {
	case config.JpegOutputFormat:
		return jpeg.Encode(file, img, nil)
	case config.PngOutputFormat:
		return png.Encode(file, img)
	case config.AvifOutputFormat:
		return avif.Encode(file, img)
	case config.WebpOutputFormat:
		return webp.Encode(file, img)
	default:
		return fmt.Errorf("unsupported output format: %v", format)
	}
}
