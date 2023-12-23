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
	Storage readers.ReaderStorage
}

func New() *ShonenJumpPlus {
	return &ShonenJumpPlus{
		Storage: readers.ReaderStorage{
			ID:     "shonenjumpplus",
			Domain: "shonenjumpplus.com",
		},
	}
}

func (s *ShonenJumpPlus) Details() readers.ReaderStorage {
	return s.Storage
}

func (s *ShonenJumpPlus) SetSession(str string) {
	s.Storage.Session = &str
}

func (s *ShonenJumpPlus) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	var c = request.Config{}
	if s.Storage.Session != nil {
		c.Cookies = []*http.Cookie{
			{
				Name:     "glsc",
				Value:    *s.Storage.Session,
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
