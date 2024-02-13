package fod

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"math"
	"net/url"
	"regexp"
	"strings"
)

func descrambleImage(img image.Image, key string) []byte {
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

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, descrambledImg, nil)

	return buf.Bytes()
}

func extractValuesFromURL(uri url.URL) (*BookCredentialsRequest, error) {
	pattern := `https://manga\.fod\.fujitv\.co\.jp/(viewer|books)/(\d+)(?:/([A-Z0-9]+))?`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(uri.String())

	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid URL format")
	}

	bookId := matches[2]

	var episodeId string
	if len(matches) == 4 {
		episodeId = matches[3]
	}

	return &BookCredentialsRequest{
		BookID:    bookId,
		EpisodeID: episodeId,
	}, nil
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
