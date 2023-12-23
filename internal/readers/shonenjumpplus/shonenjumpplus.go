package shonenjumpplus

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"image"
	"net/url"
)

type ShonenJumpPlus struct {
	Session string
}

func (s ShonenJumpPlus) Details() readers.ParserDetails {
	return readers.ParserDetails{
		ID:     "shonenjumpplus",
		Domain: "shonenjumpplus.com",
	}
}

func (s ShonenJumpPlus) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	resp, err := request.Get[ShonenJumpPlusResponse](uri.String()+".json", nil)
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
