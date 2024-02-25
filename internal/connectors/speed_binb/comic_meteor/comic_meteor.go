package comic_meteor

import (
	"net/url"

	"mary/internal/connectors/speed_binb"
	"mary/internal/static"
	"mary/pkg/request"
)

type ComicMeteor struct {
	domain string
	binb   *speed_binb.SpeedBinb
}

func (c *ComicMeteor) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusNay,
		ChapterListAvailable: false,
	}
}

func (c *ComicMeteor) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *ComicMeteor) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *ComicMeteor) Chapter(uri url.URL) (*static.Chapter, error) {
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

func (c *ComicMeteor) Pages(chapterID any, imageChan chan<- static.Image) error {
	return c.binb.Pages(chapterID.(url.URL), imageChan, nil)
}

func New() *ComicMeteor {
	domain := "comic-meteor.jp"
	return &ComicMeteor{
		domain: domain,
		binb:   speed_binb.New(domain),
	}
}
