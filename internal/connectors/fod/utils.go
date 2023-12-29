package fod

import (
	"559/internal/connectors"
	"559/pkg/request"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func processPage(uri, key string, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		img, err := request.Get[image.Image](uri, nil)
		if err != nil {
			return nil, err
		}

		return descrambleImage(img, key), nil
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}

func processOriginalPage(uri, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		return request.Get[image.Image](uri, nil)
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}

func processKeys(imageUrl string, keys []string) error {
	img, err := request.Get[image.Image](imageUrl, nil)
	if err != nil {
		return err
	}

	wi, hi := img.Bounds().Size().X, img.Bounds().Size().Y

	byteKeys := make([][]int, len(keys))

	for i, key := range keys {
		m := int(math.Floor(float64(wi / 96)))
		o := m * int(math.Floor(float64(hi/128)))

		s := make([]int, o)
		for a := 0; a < o; a++ {
			s[a] = a
		}

		s = NewRandomizer(key).Shuffle(s)

		byteKeys[i] = s
	}

	jsonData, err := json.Marshal(byteKeys)
	if err != nil {
		log.Fatal(err)
	}

	filePath := "output/keys.json"

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Data written to %s\n", filePath)
	return nil
}

/*
Options:
	:shiftRemainingEdgeDrawingPosition = false (always)
	:needsPatchForCanvasGapBug (?)
*/

func descrambleImage(img image.Image, key string) image.Image {
	wi, hi := img.Bounds().Size().X, img.Bounds().Size().Y

	descrambledImg := image.NewRGBA(image.Rect(0, 0, wi, hi))
	draw.Draw(descrambledImg, descrambledImg.Bounds(), img, image.Point{}, draw.Src)

	m := int(math.Floor(float64(wi / 96)))
	o := m * int(math.Floor(float64(hi/128)))

	s := make([]int, o)
	for a := 0; a < o; a++ {
		s[a] = a
	}

	s = NewRandomizer(key).Shuffle(s)

	for v, w := 0, len(s); v < w; v++ {
		b := +s[v]
		y := int(math.Floor(96 * float64(v%m)))
		x := 128 * int(math.Floor(float64(v/m)))
		S := int(math.Round(96 * math.Floor(float64(b%m))))
		b = int(math.Round(128 * math.Floor(float64(b/m))))

		drawRect, dstRect := image.Rect(y, x, y+96, x+128), image.Rect(S, b, S+96, b+128)

		draw.Draw(descrambledImg, dstRect, img, drawRect.Min, draw.Src)
	}

	return descrambledImg
}

func extractValuesFromURL(uri url.URL) (string, string, error) {
	pattern := `https://manga\.fod\.fujitv\.co\.jp/(viewer|books)/(\d+)(?:/([A-Z0-9]+))?`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(uri.String())

	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid URL format")
	}

	value1 := matches[2] // manga number

	var value2 string
	if len(matches) == 4 {
		value2 = matches[3]
	}

	return value1, value2, nil
}

func normalizeUrl(input string) string {
	return strings.ReplaceAll(input, "\\", "/")
}

func cleanURL(rawURL string) (string, error) {
	trimmedURL := strings.TrimSpace(rawURL)

	cleanedURL := strings.ReplaceAll(trimmedURL, "\n", "")

	_, err := url.ParseRequestURI(cleanedURL)
	if err != nil {
		return "", err
	}

	return cleanedURL, nil
}
