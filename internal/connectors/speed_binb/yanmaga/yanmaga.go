package yanmaga

import (
	"559/internal/connectors/speed_binb"
	"559/internal/static"
	"559/pkg/request"
	"net/url"
)

type Yanmanga struct {
	domain string
	binb   *speed_binb.SpeedBinb
}

func (c *Yanmanga) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusNay,
		ChapterListAvailable: true,
	}
}

func (c *Yanmanga) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *Yanmanga) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *Yanmanga) Chapter(uri url.URL) (*static.Chapter, error) {
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

func (c *Yanmanga) Pages(chapterID any, imageChan chan<- static.Image) error {
	return c.binb.Pages(chapterID.(url.URL), imageChan, nil)
}

func New() *Yanmanga {
	domain := "yanmaga.jp"
	return &Yanmanga{
		domain: domain,
		binb:   speed_binb.New(domain),
	}
}
