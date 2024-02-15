package comic_walker

import (
	"encoding/hex"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"regexp"

	"github.com/rs/zerolog/log"

	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
)

type ComicWalker struct {
	domain string
}

func New() *ComicWalker {
	return &ComicWalker{
		domain: "comic-walker.com",
	}
}

func (c *ComicWalker) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusForBookmarks,
		ChapterListAvailable: false,
	}
}

func (c *ComicWalker) ResolveType(uri url.URL) (static.UrlType, error) {
	bookRegex := regexp.MustCompile(`/contents/detail/`)
	chapterRegex := regexp.MustCompile(`/viewer/`)

	if bookRegex.MatchString(uri.Path) {
		return static.UrlTypeBook, nil
	} else if chapterRegex.MatchString(uri.Path) {
		return static.UrlTypeChapter, nil
	}

	return "", static.UnknownUrlTypeErr
}

func (c *ComicWalker) Book(uri url.URL) (*static.Book, error) {
	reqUrl := fmt.Sprintf("https://comicwalker-api.nicomanga.jp/api/v1/comicwalker/contents/%s", utils.LastURLSegment(uri.Path))
	res, err := request.Get[BookResponse](reqUrl)
	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("%s | status code: %d", reqUrl, res.Status)

	return &static.Book{
		Title:    res.Body.Data.Result.Meta.Title,
		Cover:    &res.Body.Data.Result.Meta.MainImageUrl,
		Chapters: nil,
	}, nil
}

func (c *ComicWalker) Chapter(uri url.URL) (*static.Chapter, error) {
	if !uri.Query().Has("cid") {
		return nil, fmt.Errorf("url dont have cid")
	}

	res, err := request.Get[ChapterResponse](fmt.Sprintf("https://comicwalker-api.nicomanga.jp/api/v1/comicwalker/episodes/%s", uri.Query().Get("cid")))
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    uri.Query().Get("cid"),
		Title: res.Body.Data.Result.Title,
		Error: nil,
	}, nil
}

func (c *ComicWalker) Pages(chapterID any, imageChan chan<- static.Image) error {
	response, err := request.Get[FramesResponse](fmt.Sprintf("https://comicwalker-api.nicomanga.jp/api/v1/comicwalker/episodes/%s/frames", chapterID))
	if err != nil {
		return err
	}

	indexNamer := utils.NewIndexNamer(len(response.Body.Data.Result))
	for i, page := range response.Body.Data.Result {
		if page.Meta.DrmHash != nil {
			log.Trace().Msgf("url: %s | hash: %s", page.Meta.SourceUrl, *page.Meta.DrmHash)
		} else {
			log.Trace().Msgf("url: %s | hash is nil", page.Meta.SourceUrl)
		}

		var imageFn static.ImageFn
		imageFn = func() ([]byte, error) {
			imageResponse, err := request.Get[[]byte](page.Meta.SourceUrl)
			if err != nil {
				return nil, err
			}

			if page.Meta.DrmHash != nil {
				return decodeImage(imageResponse.Body, *page.Meta.DrmHash)
			} else {
				return imageResponse.Body, nil
			}
		}

		imageChan <- static.NewImage(indexNamer.Get(i, ".jpg"), &imageFn)
	}

	close(imageChan)
	return nil
}

func decodeImage(b []byte, hash string) ([]byte, error) {
	key, err := generateKey(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %q key: %s", hash, err)
	}

	decrypted := xor(b, key)

	return decrypted, nil
}

func generateKey(t string) ([]byte, error) {
	if len(t) < 16 {
		return nil, fmt.Errorf("failed generate key")
	}

	keyBytes, err := hex.DecodeString(t[:16])
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

func xor(t []byte, e []byte) []byte {
	r, i := len(t), len(e)
	o := make([]byte, r)

	for a := 0; a < r; a++ {
		o[a] = t[a] ^ e[a%i]
	}

	return o
}
