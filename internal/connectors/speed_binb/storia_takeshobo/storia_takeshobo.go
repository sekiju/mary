package storia_takeshobo

import (
	"559/internal/connectors/speed_binb"
	"559/internal/static"
	"559/pkg/request"
	"net/url"
)

type StoriaTakeshobo struct {
	domain string
	binb   *speed_binb.SpeedBinb
}

func (c *StoriaTakeshobo) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusNay,
		ChapterListAvailable: true,
	}
}

func (c *StoriaTakeshobo) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *StoriaTakeshobo) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *StoriaTakeshobo) Chapter(uri url.URL) (*static.Chapter, error) {
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

func (c *StoriaTakeshobo) Pages(chapterID any, imageChan chan<- static.Image) error {
	return c.binb.Pages(chapterID.(url.URL), imageChan, nil)
}

func New() *StoriaTakeshobo {
	domain := "storia.takeshobo.co.jp"
	return &StoriaTakeshobo{
		domain: domain,
		binb:   speed_binb.New(domain),
	}
}
