package giga_viewer

import (
	"559/internal/config"
	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
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

func (c *GigaViewer) ResolveType(uri url.URL) (static.UrlType, error) {
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
		chapterResponse, err := request.Get[EpisodeResponse](item.Link+".json", c.withCookies())
		if err != nil {
			return nil, static.NotFoundErr
		}

		chapter := static.Chapter{
			ID:    chapterResponse.Body.ReadableProduct.Id,
			Title: chapterResponse.Body.ReadableProduct.Title,
			Error: nil,
		}

		if !chapterResponse.Body.ReadableProduct.IsPublic && !chapterResponse.Body.ReadableProduct.HasPurchased {
			chapter.Error = static.PaidChapterErr
		}

		chapters = append(chapters, chapter)
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

	chapter := static.Chapter{
		ID:    res.Body.ReadableProduct.Id,
		Title: res.Body.ReadableProduct.Title,
		Error: nil,
	}

	if !res.Body.ReadableProduct.IsPublic && !res.Body.ReadableProduct.HasPurchased {
		chapter.Error = static.PaidChapterErr
	}

	return &chapter, nil
}

func (c *GigaViewer) Pages(chapterID any, imageChan chan<- static.Image) error {
	res, err := request.Get[EpisodeResponse](fmt.Sprintf("https://%s/episode/%s.json", c.domain, chapterID), c.withCookies())
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
		processPage(page.Src, indexNamer.Get(i, ".jpg"), imageChan)
	}

	close(imageChan)
	return nil
}

func processPage(uri, fileName string, imageChan chan<- static.Image) {
	var fn static.ImageFn
	fn = func() ([]byte, error) {
		res, err := request.Get[[]byte](uri)
		if err != nil {
			return nil, err
		}

		return res.Body, err
	}

	imageChan <- static.NewImage(fileName, &fn)
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
	connectorConfig, exists := config.Data.Sites[c.domain]
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

func New(domain string) *GigaViewer {
	return &GigaViewer{domain}
}
