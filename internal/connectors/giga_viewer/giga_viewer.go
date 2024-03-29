package giga_viewer

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"

	"mary/internal/config"
	"mary/internal/static"
	"mary/internal/utils"
	"mary/pkg/request"
)

type GigaViewer struct {
	domain string
}

func (c *GigaViewer) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusOptional,
		ChapterListAvailable: true,
	}
}

func (c *GigaViewer) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeShared, nil
}

func (c *GigaViewer) Book(uri url.URL) (*static.Book, error) {
	res, err := request.Get[EpisodeResponse](uri.String()+".json", c.withCookies())
	if err != nil {
		return nil, static.NotFoundErr
	}

	rssText, err := request.Get[string](fmt.Sprintf("https://shonenjumpplus.com/rss/series/%s", res.Body.ReadableProduct.Series.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rss channel for chapters list")
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(rssText.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rss feed")
	}

	chapters := make([]static.Chapter, 0)
	for _, item := range feed.Items {
		chapterUrl, _ := url.Parse(item.Link)
		chapter, err := c.Chapter(*chapterUrl)
		if err != nil {
			return nil, err
		}

		chapters = append(chapters, *chapter)
	}

	return &static.Book{
		Title:    res.Body.ReadableProduct.Series.Title,
		Cover:    &res.Body.ReadableProduct.Series.ThumbnailUri,
		Chapters: nil,
	}, nil
}

func (c *GigaViewer) Chapter(uri url.URL) (*static.Chapter, error) {
	res, err := request.Get[EpisodeResponse](uri.String()+".json", c.withCookies())
	if err != nil {
		return nil, static.NotFoundErr
	}

	credentials, err := extractValuesFromURL(uri)
	if err != nil {
		return nil, err
	}

	chapter := static.Chapter{
		ID:    credentials,
		Title: res.Body.ReadableProduct.Title,
		Error: nil,
	}

	if !res.Body.ReadableProduct.IsPublic && !res.Body.ReadableProduct.HasPurchased {
		chapter.Error = static.PaidChapterErr
	}

	return &chapter, nil
}

func (c *GigaViewer) Pages(chapterID any, imageChan chan<- static.Image) error {
	credentials := chapterID.(*Credentials)

	res, err := request.Get[EpisodeResponse](fmt.Sprintf("https://%s/%s/%s.json", c.domain, credentials.Type, credentials.ID), c.withCookies())
	if err != nil {
		return static.NotFoundErr
	}

	if !res.Body.ReadableProduct.IsPublic && !res.Body.ReadableProduct.HasPurchased {
		return static.PaidChapterErr
	}

	pages := filterMainPages(res.Body.ReadableProduct.PageStructure.Pages)
	indexNamer := utils.NewIndexNamer(len(pages))
	for i, page := range pages {
		log.Trace().Msgf("url: %s", page.Src)

		var fn static.ImageFn
		fn = func() ([]byte, error) {
			res, err := request.Get[[]byte](page.Src)
			if err != nil {
				return nil, err
			}

			return res.Body, err
		}

		imageChan <- static.NewImage(indexNamer.Get(i, ".jpg"), &fn)
	}

	close(imageChan)
	return nil
}

func filterMainPages(pages []EpisodePage) []EpisodePage {
	filtered := make([]EpisodePage, 0)
	for _, page := range pages {
		if page.Type != "main" {
			log.Trace().Msg("skipping tshirt ad page")
			continue
		}

		filtered = append(filtered, page)
	}

	return filtered
}

func (c *GigaViewer) withCookies() request.OptsFn {
	connectorConfig, exists := config.Config.Sites[c.domain]
	return func(cf *request.Config) {
		if exists {
			log.Trace().Msgf("used cookies for %s", c.domain)
			cf.Cookies = append(cf.Cookies, &http.Cookie{
				Name:     "glsc",
				Value:    connectorConfig.Session,
				Path:     "/",
				Domain:   c.domain,
				Secure:   true,
				HttpOnly: true,
			})
		}
	}
}

func extractValuesFromURL(uri url.URL) (*Credentials, error) {
	pattern := `(magazine|episode)\/(\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(uri.String())

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid URL format")
	}

	return &Credentials{
		Type: matches[1],
		ID:   matches[2],
	}, nil
}

func New(domain string) *GigaViewer {
	return &GigaViewer{domain}
}
