package giga_viewer

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"image"
	"net/http"
	"net/url"
)

type GigaViewer struct {
	*readers.Base
}

func New(domain string) *GigaViewer {
	return &GigaViewer{Base: readers.NewBase(domain)}
}

func (g *GigaViewer) Context() *readers.Base {
	return g.Base
}

func (g *GigaViewer) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	var c = request.Config{}

	session, exists := g.Data["session"]

	if exists {
		c.Cookies = []*http.Cookie{
			{
				Name:     "glsc",
				Value:    session.(string),
				Path:     "/",
				Domain:   g.Domain,
				Secure:   true,
				HttpOnly: true,
			},
		}
	}

	resp, err := request.Get[EpisodeResponse](uri.String()+".json", &c)
	if err != nil {
		return err
	}

	fnf := utils.NewIndexNameFormatter(len(resp.ReadableProduct.PageStructure.Pages))
	for i, page := range resp.ReadableProduct.PageStructure.Pages {
		if page.Type != "main" {
			continue
		}

		// info: memorize
		src := page.Src

		var imageFunc readers.ImageFunction
		imageFunc = func() (image.Image, error) {
			return request.Get[image.Image](src, nil)
		}

		imageChan <- readers.NewReaderImage(fnf.GetName(i, ".jpg"), &imageFunc)
	}

	return nil
}
