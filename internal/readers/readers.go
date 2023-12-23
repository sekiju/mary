package readers

import (
	"image"
	"net/url"
)

type ParserDetails struct {
	ID     string
	Domain string
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

type Parser interface {
	Details() ParserDetails
	Pages(uri url.URL, imageChan chan<- ReaderImage) error
}
