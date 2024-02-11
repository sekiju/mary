package testing_utils

import (
	"image"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func getAsset(path, url string) (image.Image, error) {
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, err
		}

		return img, nil
	}

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return nil, err
	}

	file.Seek(0, io.SeekStart)
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func rootFolder(t *testing.T) func() {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	_, b, _, _ := runtime.Caller(0)
	d := filepath.Dir(path.Join(path.Dir(b), "../.."))

	err = os.Chdir(d)
	if err != nil {
		t.Fatal(err)
	}

	return func() {
		err := os.Chdir(originalWD)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func failedResultPath(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	return base + "_failed" + ext
}
