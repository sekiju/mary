package connectors

import (
	"image"
	"net/url"
)

type Base struct {
	Domain string
}

func NewBase(domain string) *Base {
	return &Base{
		Domain: domain,
	}
}

type ImageFunction func() (image.Image, error)
type ReaderImage struct {
	FileName string
	Image    ImageFunction
}

func NewConnectorImage(fn string, imf *ImageFunction) ReaderImage {
	return ReaderImage{
		FileName: fn,
		Image:    *imf,
	}
}

type Connector interface {
	Context() *Base
	Pages(uri url.URL, imageChan chan<- ReaderImage) error
}
