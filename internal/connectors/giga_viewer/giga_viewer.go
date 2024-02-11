package giga_viewer

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"image"
	"net/http"
	"net/url"
)

type GigaViewer struct {
	*connectors.Base
}

func New(domain string) *GigaViewer {
	return &GigaViewer{Base: connectors.NewBase(domain)}
}

func (g *GigaViewer) Context() *connectors.Base {
	return g.Base
}

func (g *GigaViewer) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	var c = request.Config{}

	connectorConfig, exists := config.Data.Sites[g.Domain]

	if exists {
		c.Cookies = []*http.Cookie{
			{
				Name:     "glsc",
				Value:    connectorConfig.Session,
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

		processPage(page.Src, fnf.GetName(i, ".jpg"), imageChan)
	}

	close(imageChan)

	return nil
}

func processPage(uri, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		return request.Get[image.Image](uri, nil)
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}
