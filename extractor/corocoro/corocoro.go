package corocoro

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	json "github.com/bytedance/sonic"
	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/extractor/util"
	"github.com/knst0/mdl/sdk/manga"
	"github.com/knst0/mdl/sdk/manga/pluginutil"

	"resty.dev/v3"
)

type Extractor struct {
	pluginutil.Base
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	return nil, manga.ErrChapterListingUnsupported
}

func (e *Extractor) SupportsChapterListing() bool {
	return false
}

var re = regexp.MustCompile(`https://www.corocoro.jp/chapter/(\d*)/viewer`)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 2 {
		return nil, manga.ErrInvalidChapterURL
	}

	req := resty.New().R().SetContext(ctx)
	e.ApplyCookie(req)

	res, err := req.Get(URL)
	if err != nil {
		return nil, err
	}

	html := res.String()

	chapterMainName, err := util.ExtractStringFromHTML(html, `\"chapterMainName\":\"`, `\"`)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      matches[1],
		Number:  "",
		Title:   chapterMainName,
		Index:   0,
		URL:     URL,
		MangaID: "",
	}, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	req := resty.New().R().SetContext(ctx)
	e.ApplyCookie(req)

	res, err := req.Get(chapter.URL)
	if err != nil {
		return nil, err
	}

	html := res.String()

	jsonStr, err := util.ExtractStringFromHTML(html, `,\"pages\":`, `,\"directionRightToLeft\"`)
	if err != nil {
		return nil, err
	}

	jsonStr = strings.Replace(jsonStr, `\"`, `"`, -1)

	var result pagesResult
	if err = json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}

	return pluginutil.BuildPages(len(result), ".webp", func(index int, filename string) (*manga.Page, error) {
		page := result[index]
		return &manga.Page{
			Index:    uint(index),
			URL:      page.Src,
			Filename: filename,
			Decode: func(b []byte) ([]byte, error) {
				key, err := hex.DecodeString(page.Crypto.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid key hex: %w: %w", err, manga.ErrMalformedChapterData)
				}

				iv, err := hex.DecodeString(page.Crypto.Iv)
				if err != nil {
					return nil, fmt.Errorf("invalid IV hex: %w: %w", err, manga.ErrMalformedChapterData)
				}

				block, err := aes.NewCipher(key)
				if err != nil {
					return nil, err
				}

				mode := cipher.NewCBCDecrypter(block, iv)
				mode.CryptBlocks(b, b)

				return b[:(len(b) - int(b[len(b)-1]))], nil
			},
		}, nil
	})
}

func init() {
	registry.Register("www.corocoro.jp", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
