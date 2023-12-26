package fod

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"fmt"
	"image"
	"image/draw"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Fod struct {
	ctx readers.ReaderContext
}

func New() *Fod {
	return &Fod{
		ctx: readers.ReaderContext{
			Domain: "manga.fod.fujitv.co.jp",
			Data:   map[string]any{},
		},
	}
}

func (f *Fod) Context() readers.ReaderContext {
	return f.ctx
}

func (f *Fod) UpdateData(key string, value any) {
	f.ctx.Data[key] = value
}

func (f *Fod) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	bookId, episodeId, err := extractValuesFromURL(uri)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"zk-app-version": "1.1.25",
		"zk-os-type":     "1",
		"zk-safe-search": "0",
	}

	session, sessionExists := f.ctx.Data["session"]
	if sessionExists {
		headers["zk-session-key"] = session.(string)
	}

	resp, err := request.Post[LicenceKeyResponse]("https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser", &request.Config{
		Headers: headers,
		Body: map[string]string{
			"book_id":    bookId,
			"episode_id": episodeId,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch episode: %s", err)
	}

	tryPurchaseBook, exists := f.ctx.Data["tryPurchaseBook"]
	if sessionExists && exists && tryPurchaseBook.(bool) {
		isFullVersion := strings.Contains(resp.GuardianInfoForBrowser.BookData.S3Key, "_001")
		if !isFullVersion {
			_, err := request.Post[interface{}]("https://manga.fod.fujitv.co.jp/api/purchase/buy", &request.Config{
				Headers: headers,
				Body: map[string]any{
					"buy_type": 1,
					"episodes": []map[string]any{
						{
							"episode_id":       episodeId,
							"discounted_price": 0,
							"cashback_point":   0,
						},
					},
				},
			})

			if err != nil {
				fmt.Println("Failed to purchase book")
			} else {
				fmt.Println("Purchase success")
			}

			resp, err = request.Post[LicenceKeyResponse]("https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser", &request.Config{
				Headers: headers,
				Body: map[string]string{
					"book_id":    bookId,
					"episode_id": episodeId,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to fetch episode: %s", err)
			}
		}
	}

	fmt.Printf("total pages: %d\n", resp.GuardianInfoForBrowser.BookData.PageCount)
	fmt.Println(resp.GuardianInfoForBrowser.PagesData.Keys)

	saveOriginal, exists := f.ctx.Data["saveOriginal"]

	fnf := utils.NewIndexNameFormatter(resp.GuardianInfoForBrowser.BookData.PageCount)
	for i := 1; i <= resp.GuardianInfoForBrowser.PagesData.End+1; i++ {
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

			if exists && saveOriginal.(bool) {
				return img, nil
			}

			return descrambleImage(img, resp.GuardianInfoForBrowser.PagesData.Keys[currentIndex-1]), nil
		}

		imageChan <- readers.NewReaderImage(fnf.GetName(i, ".jpg"), &imageFunc)
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
