package test_utils

import (
	"559/internal/connectors"
	"github.com/corona10/goimagehash"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func ReaderTest(t *testing.T, r connectors.Connector, assetPath, assetUrl, chapterUrl string) {
	restore := RootFolder(t)
	defer restore()

	expected, err := GetAsset(assetPath, assetUrl)
	if err != nil {
		t.Fatal(err)
	}

	uri, err := url.Parse(chapterUrl)
	if err != nil {
		t.Fatal(err)
	}

	imageChan := make(chan connectors.ReaderImage)

	go func() {
		defer close(imageChan)

		err = r.Pages(*uri, imageChan)
		if err != nil {
			t.Fatal(err)
		}
	}()

	var img image.Image
	for ic := range imageChan {
		img, err = ic.Image()
		if err != nil {
			t.Fatalf("failed to get image: %s", err)
		}
		break
	}

	hashExpected, _ := goimagehash.AverageHash(expected)
	hashResult, _ := goimagehash.AverageHash(img)

	distance, _ := hashResult.Distance(hashExpected)
	if distance > 1 {
		saveImage(img, addFailedSuffix(assetPath))
		t.Fatalf("Image hashes do not match. Distance: %d", distance)
	}
}

func addFailedSuffix(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)

	return base + "_failed" + ext
}

func GetAsset(path, url string) (image.Image, error) {
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

func RootFolder(t *testing.T) func() {
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

func saveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
