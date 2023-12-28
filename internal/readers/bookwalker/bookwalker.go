package bookwalker

import (
	"559/internal/readers"
	"net/url"
)

// https://github.com/manga-download/hakuneko/blob/master/src/web/mjs/connectors/templates/Publus.mjs

/*

 */

type BookWalker struct {
	*readers.Base
}

func New() *BookWalker {
	return &BookWalker{
		Base: readers.NewBase("bookwalker.jp"),
	}
}

func (b *BookWalker) Context() *readers.Base {
	return b.Base
}

func (b *BookWalker) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	return nil
}
