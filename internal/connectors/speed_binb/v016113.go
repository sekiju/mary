package speed_binb

import (
	"559/internal/static"
	"net/url"
)

func handleV016113(uri url.URL, apiUrl string, imageChan chan<- static.Image) error {
	return handleV016130(uri, apiUrl, imageChan)
}
