package speed_binb

import (
	"559/internal/connectors"
	"net/url"
)

func handleV016113(uri url.URL, apiUrl string, imageChan chan<- connectors.ReaderImage) error {
	return handleV016130(uri, apiUrl, imageChan)
}
