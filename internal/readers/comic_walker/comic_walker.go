package comic_walker

import (
	"559/internal/readers"
	"559/internal/utils"
	"559/internal/utils/request"
	"bytes"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
)

type ComicWalker struct {
	ctx readers.ReaderContext
}

func New() *ComicWalker {
	return &ComicWalker{
		ctx: readers.ReaderContext{
			Domain: "comic-walker.com",
			Data:   map[string]any{},
		},
	}
}

func (c *ComicWalker) Context() readers.ReaderContext {
	return c.ctx
}

func (c *ComicWalker) UpdateData(key string, value any) {
	c.ctx.Data[key] = value
}

func (c *ComicWalker) Pages(uri url.URL, imageChan chan<- readers.ReaderImage) error {
	// todo: session support

	if !uri.Query().Has("cid") {
		return fmt.Errorf("url dont have cid")
	}

	resp, err := request.Get[FramesResponse](fmt.Sprintf("https://comicwalker-api.nicomanga.jp/api/v1/comicwalker/episodes/%s/frames", uri.Query().Get("cid")), nil)
	if err != nil {
		return err
	}

	fnf := utils.NewIndexNameFormatter(len(resp.Data.Result))
	for i, page := range resp.Data.Result {

		// info: memorize
		src, hash := page.Meta.SourceUrl, page.Meta.DrmHash

		var imageFunc readers.ImageFunction
		imageFunc = func() (image.Image, error) {
			img, err := request.Get[[]byte](src, nil)
			if err != nil {
				return nil, err
			}

			return decodeImage(img, hash)
		}

		imageChan <- readers.NewReaderImage(fnf.GetName(i, ".jpg"), &imageFunc)
	}

	return nil
}

func decodeImage(b []byte, hash string) (image.Image, error) {
	key, err := generateKey(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to generate %q key: %s", hash, err)
	}

	decrypted := xor(b, key)

	img, err := jpeg.Decode(bytes.NewReader(decrypted))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func generateKey(t string) ([]byte, error) {
	if len(t) < 16 {
		return nil, fmt.Errorf("failed generate key")
	}

	keyBytes, err := hex.DecodeString(t[:16])
	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

func xor(t []byte, e []byte) []byte {
	r, i := len(t), len(e)
	o := make([]byte, r)

	for a := 0; a < r; a++ {
		o[a] = t[a] ^ e[a%i]
	}

	return o
}
