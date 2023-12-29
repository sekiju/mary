package bookwalker

import (
	"559/internal/connectors"
	"net/url"
)

// https://github.com/manga-download/hakuneko/blob/master/src/web/mjs/connectors/templates/Publus.mjs

/*

 */

type BookWalker struct {
	*connectors.Base
}

func New() *BookWalker {
	return &BookWalker{
		Base: connectors.NewBase("bookwalker.jp"),
	}
}

func (b *BookWalker) Context() *connectors.Base {
	return b.Base
}

func (b *BookWalker) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	close(imageChan)

	return nil
}

func processPage(uri, fileName string, imageChan chan<- connectors.ReaderImage) {}
