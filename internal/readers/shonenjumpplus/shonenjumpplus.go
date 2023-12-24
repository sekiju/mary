package shonenjumpplus

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"image"
	"net/http"
	"net/url"
)

type ShonenJumpPlus struct {
	ctx readers.ReaderContext
}

func New() *ShonenJumpPlus {
	return &ShonenJumpPlus{
		ctx: readers.ReaderContext{
			Domain: "shonenjumpplus.com",
			Data:   map[string]any{},
		},
	}
}

func (s *ShonenJumpPlus) Context() readers.ReaderContext {
	return s.ctx
}

func (s *ShonenJumpPlus) UpdateData(k string, v any) {
	s.ctx.Data[k] = v
}

func (s *ShonenJumpPlus) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	var c = request.Config{}

	session, exists := s.ctx.Data["session"]

	if exists {
		c.Cookies = []*http.Cookie{
			{
				Name:     "glsc",
				Value:    session.(string),
				Path:     "/",
				Domain:   "shonenjumpplus.com",
				Secure:   true,
				HttpOnly: true,
			},
		}
	}

	resp, err := request.Get[ShonenJumpPlusResponse](uri.String()+".json", &c)
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
