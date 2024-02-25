package comic_walker

import (
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	_ "image/jpeg"
	_ "image/png"
	"mary/internal/static"
	"mary/internal/utils"
	"mary/pkg/request"
	"net/url"
	"regexp"
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

func (c *ComicWalker) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func exportCodes(uri url.URL) (*ExportedCodes, error) {
	reg := regexp.MustCompile("KC_[0-9]+(?:_E|_S)?")

	matches := reg.FindAllString(uri.Path, -1)
	if len(matches) != 2 {
		return nil, fmt.Errorf("invalid URL format")
	}

	return &ExportedCodes{
		Work:    matches[0],
		Episode: matches[1],
	}, nil
}

func (c *ComicWalker) Book(uri url.URL) (*static.Book, error) {
	codes, err := exportCodes(uri)
	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("https://comic-walker.com/api/contents/details/work?workCode=%s", codes.Work)
	log.Trace().Msg(reqUrl)

	res, err := request.Get[WorkResponse](reqUrl)
	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("%s | status code: %d", reqUrl, res.Status)

	return &static.Book{
		Title:    res.Body.Work.Title,
		Cover:    &res.Body.Work.BookCover,
		Chapters: nil,
	}, nil
}

func (c *ComicWalker) Chapter(uri url.URL) (*static.Chapter, error) {
	codes, err := exportCodes(uri)
	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeCode=%s&episodeType=latest", codes.Work, codes.Episode)
	log.Trace().Msg(reqUrl)

	res, err := request.Get[EpisodeResponse](reqUrl)
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    res.Body.Episode.Id,
		Title: res.Body.Episode.Title,
		Error: nil,
	}, nil
}

func (c *ComicWalker) Pages(chapterID any, imageChan chan<- static.Image) error {
	reqUrl := fmt.Sprintf("https://comic-walker.com/api/contents/viewer?episodeId=%s&imageSizeType=width%s1284", chapterID, "%3A")
	log.Trace().Msg(reqUrl)

	res, err := request.Get[ViewerResponse](reqUrl)
	if err != nil {
		return err
	}

	indexNamer := utils.NewIndexNamer(len(res.Body.Manuscripts))
	for i, page := range res.Body.Manuscripts {
		var imageFn static.ImageFn
		imageFn = func() ([]byte, error) {
			imageResponse, err := request.Get[[]byte](page.DrmImageUrl)
			if err != nil {
				return nil, err
			}

			return decodeImage(imageResponse.Body, page.DrmHash)
		}

		imageChan <- static.NewImage(indexNamer.Get(i, ".webp"), &imageFn)
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
