package fod

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Fod struct {
	Session   string
	DebugKeys bool
}

func (f Fod) Details() readers.ParserDetails {
	return readers.ParserDetails{
		ID:     "fod",
		Domain: "manga.fod.fujitv.co.jp",
	}
}

func (f Fod) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	bookId, episodeId, err := extractValuesFromURL(uri)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"zk-app-version": "1.1.24",
		"zk-os-type":     "1",
		"zk-safe-search": "0",
	}

	if len(f.Session) > 0 {
		headers["zk-session-key"] = f.Session
	}

	resp, err := request.Post[LicenceKeyResponse]("https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser", &request.Config{
		Headers: headers,
		Body: map[string]string{
			"book_id":    bookId,
			"episode_id": episodeId,
		},
	})
	if err != nil {
		return err
	}

	fnf := utils.NewIndexNameFormatter(resp.GuardianInfoForBrowser.BookData.PageCount)
	for i := 1; i <= resp.GuardianInfoForBrowser.BookData.PageCount; i++ {
		// info: memorize index
		currentIndex := i

		var imageFunc readers.ImageFunction
		imageFunc = func() (image.Image, error) {
			imageUrl, err := cleanURL(resp.GuardianInfoForBrowser.GUARDIANSERVER + normalizeUrl(resp.GuardianInfoForBrowser.BookData.S3Key) + strconv.Itoa(currentIndex) + ".jpg?" + resp.GuardianInfoForBrowser.ADDITIONALQUERYSTRING)
			if err != nil {
				return nil, err
			}

			img, err := request.Get[image.Image](imageUrl, nil)
			if err != nil {
				return nil, err
			}

			return descrambleImage(img, resp.GuardianInfoForBrowser.PagesData.Keys[currentIndex-1]), nil
		}

		imageChan <- readers.NewReaderImage(fnf.GetName(i, ".jpg"), &imageFunc)
	}

	if f.DebugKeys {
		keys, _ := json.MarshalIndent(resp.GuardianInfoForBrowser.PagesData.Keys, "", " ")
		_ = ioutil.WriteFile("./tests/images-keys.json", keys, 0644)
	}

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
	pattern := `https://manga\.fod\.fujitv\.co\.jp/viewer/(\d+)/([A-Z0-9]+/)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(uri.String())

	if len(matches) < 3 {
		return "", "", fmt.Errorf("invalid URL format")
	}

	value1 := matches[1]
	value2 := matches[2]

	value2 = strings.ReplaceAll(value2, "/", "")

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
