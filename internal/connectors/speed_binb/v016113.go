package speed_binb

import (
	"github.com/sekiju/rq"
	"mary/internal/static"
	"net/url"
)

func handleV016113(uri url.URL, apiUrl string, requestOpt rq.OptsFn, imageChan chan<- static.Image) error {
	return handleV016130(uri, apiUrl, requestOpt, imageChan)
}
