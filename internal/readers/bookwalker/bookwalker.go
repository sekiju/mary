package bookwalker

import (
	"559/internal/readers"
	"net/url"
)

// https://github.com/manga-download/hakuneko/blob/master/src/web/mjs/connectors/templates/Publus.mjs

/*

 */

type BookWalker struct {
	ctx readers.ReaderContext
}

func New() *BookWalker {
	return &BookWalker{
		ctx: readers.ReaderContext{
			Domain: "bookwalker.jp",
			Data:   map[string]any{},
		},
	}
}

func (c *BookWalker) Context() readers.ReaderContext {
	return c.ctx
}

func (c *BookWalker) UpdateData(key string, value any) {
	c.ctx.Data[key] = value
}

func (c *BookWalker) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	return nil
}
