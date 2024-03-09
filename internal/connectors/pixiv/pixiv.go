package pixiv

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"mary/internal/config"
	"mary/internal/static"
	"mary/internal/utils"
	"mary/pkg/request"
)

type Pixiv struct {
	domain string
}

func (c *Pixiv) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusOptional,
		ChapterListAvailable: false,
	}
}

func (c *Pixiv) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *Pixiv) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *Pixiv) Chapter(uri url.URL) (*static.Chapter, error) {
	id := utils.LastURLSegment(uri.Path)

	illustUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s", id)
	log.Trace().Msgf("illustUrl: %s", illustUrl)

	res, err := request.Get[IllustResponse](illustUrl, c.withCookies())
	if err != nil {
		return nil, fmt.Errorf("failed to access episode api: %s", err)
	}

	if res.Body.Error {
		return nil, fmt.Errorf("error status")
	}

	chapter := static.Chapter{
		ID:    id,
		Title: res.Body.Body.Title,
		Error: nil,
	}

	if res.Body.Body.NoLoginData != nil {
		for _, tag := range res.Body.Body.Tags.Tags {
			if tag.Tag == "R-18" {
				chapter.Error = static.LoginRequiredErr
			}
		}
	}

	return &chapter, nil
}

func (c *Pixiv) Pages(chapterID any, imageChan chan<- static.Image) error {
	pagesUrl := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/pages", chapterID)
	log.Trace().Msgf("pagesUrl: %s", pagesUrl)

	res, err := request.Get[PagesResponse](pagesUrl, c.withCookies())
	if err != nil {
		return fmt.Errorf("failed to access episode api: %s", err)
	}

	if res.Body.Error {
		return fmt.Errorf("error status")
	}

	indexNamer := utils.NewIndexNamer(len(res.Body.Body))
	for i, page := range res.Body.Body {
		log.Trace().Msgf("url: %s", page.Urls.Original)

		ext := filepath.Ext(path.Base(page.Urls.Original))

		var fn static.ImageFn
		fn = func() ([]byte, error) {
			res, err := request.Get[[]byte](page.Urls.Original, request.SetHeader("Referer", "https://pixiv.net/"))
			if err != nil {
				return nil, err
			}

			return res.Body, nil
		}

		imageChan <- static.NewImage(indexNamer.Get(i, ext), &fn)
	}

	close(imageChan)
	return nil
}

func (c *Pixiv) withCookies() request.OptsFn {
	connectorConfig, exists := config.Config.Sites[c.domain]
	return func(cf *request.Config) {
		if exists {
			log.Trace().Msg("used cookies for .pixiv.net")
			cf.Cookies = append(cf.Cookies, &http.Cookie{
				Name:     "PHPSESSID",
				Value:    connectorConfig.Session,
				Path:     "/",
				Domain:   ".pixiv.net",
				Secure:   true,
				HttpOnly: true,
			})
		}
	}
}

func New() *Pixiv {
	return &Pixiv{
		domain: "www.pixiv.net",
	}
}
