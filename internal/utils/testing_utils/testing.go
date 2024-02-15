package testing_utils

import (
	"bytes"
	"fmt"
	"image"
	"net/url"
	"strings"
	"testing"

	"github.com/corona10/goimagehash"
	"github.com/rs/zerolog/log"

	"559/internal/static"
	"559/internal/utils"
)

func Connector(t *testing.T, c static.Connector, assetUrl, chapterUrl string) {
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

	imageChan := make(chan static.Image)
	go func() {
		defer close(imageChan)

		err = c.Pages(chapter.ID, imageChan)
		if err != nil {
			t.Fatal(err)
		}
	}()

	var b []byte
	for ic := range imageChan {
		b, err = ic.ImageFn()
		if err != nil {
			t.Fatalf("failed to get image: %s", err)
		}

		break
	}

	reader := bytes.NewReader(b)
	img, format, err := image.Decode(reader)
	if err != nil {
		t.Fatalf("failed to decode image: %s", err)
	}

	connectorName := strings.Replace(c.Data().Domain, ".", "_", -1)
	assetPath := fmt.Sprintf("tests/assets/%s.%s", connectorName, format)

	expected, err := getAsset(assetPath, assetUrl)
	if err != nil {
		log.Trace().Msgf("path: %s | url: %s", assetPath, assetUrl)
		t.Fatal(err)
	}

	hashExpected, _ := goimagehash.AverageHash(expected)
	hashResult, _ := goimagehash.AverageHash(img)

	distance, _ := hashResult.Distance(hashExpected)
	if distance > 1 {
		err = utils.WriteImageBytes(failedResultPath(assetPath), b)
		if err != nil {
			t.Error(err)
		}

		t.Fatalf("Image hashes do not match. Distance: %d", distance)
	}
}
