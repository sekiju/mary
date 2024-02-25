package comic_valkyrie

import (
	"mary/internal/connectors/speed_binb"
	"mary/internal/static"
	"mary/pkg/request"
	"net/url"
)

type ComicValkyrie struct {
	domain string
	binb   *speed_binb.SpeedBinb
}

func (c *ComicValkyrie) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusNay,
		ChapterListAvailable: true,
	}
}

func (c *ComicValkyrie) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *ComicValkyrie) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *ComicValkyrie) Chapter(uri url.URL) (*static.Chapter, error) {
	document, err := request.Document(uri.String())
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    uri,
		Title: document.Find("title").Text(),
		Error: nil,
	}, nil
}

func (c *ComicValkyrie) Pages(chapterID any, imageChan chan<- static.Image) error {
	return c.binb.Pages(chapterID.(url.URL), imageChan, nil)
}

func New() *ComicValkyrie {
	domain := "www.comic-valkyrie.com"
	return &ComicValkyrie{
		domain: domain,
		binb:   speed_binb.New(domain),
	}
}
