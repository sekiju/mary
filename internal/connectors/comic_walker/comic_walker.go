package comic_walker

import (
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
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
	*connectors.Base
}

func New() *ComicWalker {
	return &ComicWalker{
		Base: connectors.NewBase("comic-walker.com"),
	}
}

func (c *ComicWalker) Context() *connectors.Base {
	return c.Base
}

func (c *ComicWalker) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
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
		processPage(page.Meta.SourceUrl, page.Meta.DrmHash, fnf.GetName(i, ".jpg"), imageChan)
	}

	close(imageChan)

	return nil
}

func processPage(uri, hash string, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		img, err := request.Get[[]byte](uri, nil)
		if err != nil {
			return nil, err
		}

		return decodeImage(img, hash)
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
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
