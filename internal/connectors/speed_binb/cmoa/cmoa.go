package cmoa

import (
	"559/internal/config"
	"559/internal/connectors/speed_binb"
	"559/internal/static"
	"559/pkg/request"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Cmoa struct {
	domain string
	binb   *speed_binb.SpeedBinb
}

func (c *Cmoa) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusRequired,
		ChapterListAvailable: true,
	}
}

func (c *Cmoa) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *Cmoa) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *Cmoa) Chapter(uri url.URL) (*static.Chapter, error) {
	document, err := request.Document(uri.String(), c.withCookies())
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    uri,
		Title: document.Find("title").Text(),
		Error: nil,
	}, nil
}

func (c *Cmoa) Pages(chapterID any, imageChan chan<- static.Image) error {
	opts := c.withCookies()
	return c.binb.Pages(chapterID.(url.URL), imageChan, &opts)
}

func (c *Cmoa) withCookies() request.OptsFn {
	connectorConfig, exists := config.Data.Sites[c.domain]
	return func(cf *request.Config) {
		if exists {
			log.Trace().Msgf("used cookies for %s", c.domain)
			cf.Headers["Cookie"] = connectorConfig.Session
		}
	}
}

func New() *Cmoa {
	domain := "www.cmoa.jp"
	return &Cmoa{
		domain: domain,
		binb:   speed_binb.New(domain),
	}
}
