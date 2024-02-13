package speed_binb

import (
	"559/internal/static"
	"559/pkg/request"
	"net/url"
	"strings"
)

type SpeedBinb struct {
	domain string
}

func (c *SpeedBinb) Pages(uri url.URL, imageChan chan<- static.Image, requestOpts *request.OptsFn) error {
	doc, err := request.Document(uri.String(), unwrapRequestOpts(requestOpts))
	if err != nil {
		return err
	}

	pagesContent := doc.Find("div#content.pages").First()
	ptbinb, ptbinbExists := pagesContent.Attr("data-ptbinb")
	ptbinbCid, ptbinbCidExists := pagesContent.Attr("data-ptbinb-cid")
	q := uri.Query()

	if ptbinbExists && ptbinbCidExists {
		q.Set("cid", ptbinbCid)
		uri.RawQuery = q.Encode()
		return handleV016113(uri, ptbinb, imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") && q.Has("u0") && q.Has("u1") {
		return handleV016452(uri, ptbinb, unwrapRequestOpts(requestOpts), imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") && q.Has("u1") {
		return handleV016201(uri, imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") {
		return handleV016130(uri, ptbinb, imageChan)
	} else {
		return handleV016061(uri, imageChan, pagesContent)
	}
}

func unwrapRequestOpts(opts *request.OptsFn) request.OptsFn {
	if opts == nil {
		return func(config *request.Config) {}
	}

	return *opts
}

func New(domain string) *SpeedBinb {
	return &SpeedBinb{domain: domain}
}
