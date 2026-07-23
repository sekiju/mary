package corocoro

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/sdk/manga"
	"regexp"
	"resty.dev/v3"
	"strings"
)

type Extractor struct {
	settings *manga.Settings
}

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	return nil, manga.ErrChapterListingUnsupported
}

func (e *Extractor) SupportsChapterListing() bool {
	return false
}

var re = regexp.MustCompile(`https://www.corocoro.jp/chapter/(\d*)/viewer`)

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 2 {
		return nil, manga.ErrInvalidChapterURL
	}

	req := resty.New().R()
	if e.settings.Cookie != nil {
		req.SetHeader("Cookie", *e.settings.Cookie)
	}

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

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	req := resty.New().R()
	if e.settings.Cookie != nil {
		req.SetHeader("Cookie", *e.settings.Cookie)
	}

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

	pages := make([]*manga.Page, len(result))
	padRenamer := renamer.New(len(result))

	for index, page := range result {
		pages[index] = &manga.Page{
			Index:    uint(index),
			URL:      page.Src,
			Filename: padRenamer.Name(index, ".webp"),
			Decode: func(b []byte) ([]byte, error) {
				key, err := hex.DecodeString(page.Crypto.Key)
				if err != nil {
					return nil, fmt.Errorf("invalid key hex: %w", err)
				}

				iv, err := hex.DecodeString(page.Crypto.Iv)
				if err != nil {
					return nil, fmt.Errorf("invalid IV hex: %w", err)
				}

				block, err := aes.NewCipher(key)
				if err != nil {
					return nil, err
				}

				mode := cipher.NewCBCDecrypter(block, iv)
				mode.CryptBlocks(b, b)

				return b[:(len(b) - int(b[len(b)-1]))], nil
			},
		}
	}

	return pages, nil
}

func (e *Extractor) SetSettings(settings manga.Settings) {
	e.settings = &settings
}

func New() (manga.Extractor, error) {
	return &Extractor{settings: &manga.Settings{}}, nil
}
