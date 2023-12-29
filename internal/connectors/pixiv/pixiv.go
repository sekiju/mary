package pixiv

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"path"
)

type Pixiv struct {
	*connectors.Base
}

func New() *Pixiv {
	return &Pixiv{
		Base: connectors.NewBase("www.pixiv.net"),
	}
}

func (p *Pixiv) Context() *connectors.Base {
	return p.Base
}

func (p *Pixiv) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	lastPart := path.Base(uri.Path)

	var httpConfig request.Config
	connectorConfig, exists := config.State.Sites[p.Domain]

	if exists {
		httpConfig.Cookies = []*http.Cookie{
			{
				Name:     "PHPSESSID",
				Value:    connectorConfig.Session,
				Path:     "/",
				Domain:   ".pixiv.net",
				Secure:   true,
				HttpOnly: true,
			},
		}
	}

	resp, err := request.Get[PagesResponse](fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/pages", lastPart), &httpConfig)
	if err != nil {
		return fmt.Errorf("failed to access episode api: %s", err)
	}

	fnf := utils.NewIndexNameFormatter(len(resp.Body))
	for i, page := range resp.Body {
		processPage(page.Urls.Original, fnf.GetName(i, ".jpg"), imageChan)
	}

	close(imageChan)

	return nil
}

func processPage(uri string, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		return request.Get[image.Image](uri, &request.Config{Headers: map[string]string{
			"Referer": "https://pixiv.net/",
		}})
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}
