package test

import (
	"crypto/md5"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"mary/internal/static"
	"mary/internal/utils"
	"mary/pkg/request"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Connector(t *testing.T, c static.Connector, assetUrl, chapterUrl string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	skipTests := os.Getenv("IGNORE_TESTS")
	if skipTests != "" {
		skipList := strings.Split(skipTests, ",")
		for _, skipTest := range skipList {
			if strings.Contains(skipTest, c.Data().Domain) {
				t.Skipf("skipping tests for %q", c.Data().Domain)
			}
		}
	}

	restoreFolder := rootFolder(t)
	defer restoreFolder()

	uri, err := url.Parse(chapterUrl)
	if err != nil {
		t.Error(err)
	}

	urlType, err := c.ResolveType(*uri)
	if err != nil {
		t.Error(err)
	}

	if urlType == "BOOK" {
		t.Fatal("book urls unsupported, use chapter viewer url")
	}

	chapter, err := c.Chapter(*uri)
	if err != nil {
		t.Error(err)
	}

	if chapter.Error != nil {
		t.Error(chapter.Error)
	}

	downloadedImg, downloadedFormat := singlePage(t, c, chapter.ID)

	connectorName := strings.Replace(c.Data().Domain, ".", "_", -1)
	expectedPath := fmt.Sprintf("test/assets/%s%s", connectorName, downloadedFormat)

	expectedImg, err := getAsset(expectedPath, assetUrl)
	if err != nil {
		t.Fatal(err)
	}

	if !compareBytes(downloadedImg, expectedImg) {
		ext := filepath.Ext(expectedPath)
		expectedFailedPath := strings.TrimSuffix(expectedPath, ext) + "_failed" + ext

		if err = utils.WriteImageBytes(expectedFailedPath, downloadedImg); err != nil {
			t.Error(err)
		}

		t.Fatal("images hashes do not match")
	}
}

func compareBytes(downloaded, expected []byte) bool {
	downloadedHash := md5.Sum(downloaded)
	expectedHash := md5.Sum(expected)

	return downloadedHash == expectedHash
}

func getAsset(path, uri string) ([]byte, error) {
	if _, err := os.Stat(path); err == nil {
		return utils.ReadImageBytes(path)
	}

	img, err := request.Get[[]byte](uri)
	if err != nil {
		return nil, err
	}

	err = utils.WriteImageBytes(path, img.Body)
	if err != nil {
		return nil, err
	}

	return img.Body, nil
}

func singlePage(t *testing.T, c static.Connector, chapterId any) ([]byte, string) {
	imageChan := make(chan static.Image)
	go func() {
		defer close(imageChan)

		if err := c.Pages(chapterId, imageChan); err != nil {
			t.Error(err)
		}
	}()

	for ic := range imageChan {
		img, err := ic.ImageFn()
		if err != nil {
			t.Errorf("failed to get image: %s", err)
		}

		return img, path.Ext(ic.FileName)
	}

	return nil, ""
}

func rootFolder(t *testing.T) func() {
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	_, b, _, _ := runtime.Caller(0)
	d := filepath.Dir(path.Join(b, ".."))

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
