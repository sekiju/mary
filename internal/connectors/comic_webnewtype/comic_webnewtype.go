package comic_webnewtype

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"

	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
)

type ComicWebNewtype struct {
	domain string
}

func (c *ComicWebNewtype) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusForBookmarks,
		ChapterListAvailable: true,
	}
}

func (c *ComicWebNewtype) ResolveType(uri url.URL) (static.UrlType, error) {
	parts := strings.Split(strings.Trim(uri.Path, "/"), "/")

	if len(parts) == 2 && parts[0] == "contents" {
		return static.UrlTypeBook, nil
	} else if len(parts) == 3 && parts[0] == "contents" {
		return static.UrlTypeChapter, nil
	}

	return "", static.UnknownUrlTypeErr
}

func (c *ComicWebNewtype) Book(uri url.URL) (*static.Book, error) {
	doc, err := request.Document(uri.String())
	if err != nil {
		return nil, err
	}

	book := static.Book{
		Title:    doc.Find("h1.contents__ttl").Text(),
		Cover:    nil,
		Chapters: make([]static.Chapter, 0),
	}

	img, exists := doc.Find("div.contents__thumb-comic > img").Attr("src")
	if exists {
		book.Cover = &img
	}

	more := true
	page := 1
	for more {
		uri := fmt.Sprintf(utils.JoinURL(uri.String(), fmt.Sprintf("more/%d/", page)))
		res, err := request.Get[ContentsMoreResponse](uri)
		if err != nil {
			return nil, err
		}

		doc, err = goquery.NewDocumentFromReader(strings.NewReader(res.Body.Html))
		if err != nil {
			return nil, err
		}

		selector := doc.Find("li a h2.detail__txt--ttl-sub")
		if selector.Size() == 0 {
			more = false
		}

		selector.Each(func(i int, selection *goquery.Selection) {
			uri, exists = selection.Closest("a").Attr("href")
			if !exists {
				log.Trace().Msg("anchor without href")
			}

			book.Chapters = append(book.Chapters, static.Chapter{
				ID:    strings.TrimLeft(uri, "contents"),
				Title: selection.Text(),
				Error: nil,
			})
		})

		page++
	}

	return &book, nil
}

func (c *ComicWebNewtype) Chapter(uri url.URL) (*static.Chapter, error) {
	id := strings.TrimLeft(uri.Path, "contents/")
	log.Trace().Msgf("id: %s", id)

	doc, err := request.Document(uri.String())
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    id,
		Title: doc.Find("h2.contents__ttl--comic").Text(),
		Error: nil,
	}, nil
}

func (c *ComicWebNewtype) Pages(chapterID any, imageChan chan<- static.Image) error {
	uri := fmt.Sprintf("https://comic.webnewtype.com/contents/%s/json/", chapterID)
	log.Trace().Msgf("pages url: %s", uri)

	res, err := request.Get[[]string](uri)
	if err != nil {
		return err
	}

	indexNamer := utils.NewIndexNamer(len(res.Body))
	for i, page := range res.Body {

		imgUri := utils.JoinURL("https://comic.webnewtype.com", page)
		index := strings.Index(imgUri, "/h")
		if index != -1 {
			imgUri = imgUri[:index]
		}

		ext := filepath.Ext(path.Base(imgUri))

		log.Trace().Msgf("url: %s", imgUri)

		var imageFn static.ImageFn
		imageFn = func() ([]byte, error) {
			imageResponse, err := request.Get[[]byte](imgUri)
			if err != nil {
				return nil, err
			}

			return imageResponse.Body, nil
		}

		imageChan <- static.NewImage(indexNamer.Get(i, ext), &imageFn)
	}

	close(imageChan)
	return nil
}

func New() *ComicWebNewtype {
	return &ComicWebNewtype{domain: "comic.webnewtype.com"}
}
