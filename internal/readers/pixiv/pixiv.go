package pixiv

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"path"
)

type Pixiv struct {
	ctx readers.ReaderContext
}

func (p *Pixiv) Context() readers.ReaderContext {
	return p.ctx
}

func (p *Pixiv) UpdateData(k string, v any) {
	p.ctx.Data[k] = v
}

func (p *Pixiv) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	lastPart := path.Base(uri.Path)

	var httpConfig request.Config
	session, exists := p.ctx.Data["session"]
	if exists {
		httpConfig.Cookies = []*http.Cookie{
			{
				Name:     "PHPSESSID",
				Value:    session.(string),
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
		memoPage := page

		var imageFunc readers.ImageFunction
		imageFunc = func() (image.Image, error) {
			return request.Get[image.Image](memoPage.Urls.Original, &request.Config{Headers: map[string]string{
				"Referer": "https://pixiv.net/",
			}})
		}

		imageChan <- readers.NewReaderImage(fnf.GetName(i, ".jpg"), &imageFunc)
	}

	return nil
}

func New() *Pixiv {
	return &Pixiv{
		ctx: readers.ReaderContext{
			Domain: "www.pixiv.net",
			Data:   map[string]any{},
		},
	}
}
