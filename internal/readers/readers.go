package readers

import (
	"image"
	"net/url"
)

type ReaderStorage struct {
	ID      string
	Domain  string
	Session *string
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
	Details() ReaderStorage
	SetSession(str string)

	Pages(uri url.URL, imageChan chan<- ReaderImage) error
}
