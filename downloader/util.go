package downloader

import (
	"bytes"
	"context"
	"fmt"
	"github.com/knst0/mdl/sdk/manga"
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
