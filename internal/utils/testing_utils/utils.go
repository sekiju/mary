package testing_utils

import (
	"fmt"
	"image"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"559/internal/utils"
	"559/pkg/request"
)

func getAsset(path, uri string) (image.Image, error) {
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("err: %v", err)
		}

		return img, nil
	}

	img, err := request.Get[image.Image](uri)
	if err != nil {
		return nil, err
	}

	err = utils.WriteImage(path, img.Body)
	if err != nil {
		return nil, err
	}

	return img.Body, nil
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

func BenchmarkRootFolder(t *testing.B) func() {
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
