package newtype

import (
	"559/internal/connectors"
	"559/internal/utils"
	request "559/pkg/request"
	"image"
	"net/url"
	"strings"
)

type Newtype struct {
	*connectors.Base
}

func New() *Newtype {
	return &Newtype{Base: connectors.NewBase("comic.webnewtype.com")}
}

func (n *Newtype) Context() *connectors.Base {
	return n.Base
}

func (n *Newtype) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	resp, err := request.Get[[]string](request.JoinURL(uri.String(), "json"), nil)
	if err != nil {
		return err
	}

	fnf := utils.NewIndexNameFormatter(len(resp))
	for i, src := range resp {
		processPage(src, fnf.GetName(i, ".jpg"), imageChan)
	}

	return nil
}

func processPage(uri, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		parts := strings.Split(uri, "/h1200")
		return request.Get[image.Image](request.JoinURL("https://comic.webnewtype.com", parts[0]), nil)
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}
