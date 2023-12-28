package readers

import (
	"image"
	"net/url"
)

type Base struct {
	Domain string
	Data   map[string]any
}

func NewBase(domain string) *Base {
	return &Base{
		Domain: domain,
		Data:   map[string]any{},
	}
}

type ImageFunction func() (image.Image, error)
type ReaderImage struct {
	FileName string
	Image    ImageFunction
}

func NewReaderImage(fn string, imf *ImageFunction) ReaderImage {
	return ReaderImage{
		FileName: fn,
		Image:    *imf,
	}
}

type Reader interface {
	Context() *Base
	Pages(uri url.URL, imageChan chan<- ReaderImage) error
}
