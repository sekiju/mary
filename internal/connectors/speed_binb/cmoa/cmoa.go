package cmoa

import (
	"github.com/rs/zerolog/log"
	"github.com/sekiju/rq"
	"mary/internal/config"
	"mary/internal/connectors/speed_binb"
	"mary/internal/static"
	"mary/internal/utils"
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
	document, err := utils.Document(uri.String(), c.withCookies())
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

func (c *Cmoa) withCookies() rq.OptsFn {
	connectorConfig, exists := config.Config.Sites[c.domain]
	return func(cf *rq.Opts) {
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
