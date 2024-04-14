package speed_binb

import (
	"github.com/sekiju/rq"
	"mary/internal/static"
	"mary/internal/utils"
	"net/url"
	"strings"
)

type SpeedBinb struct {
	domain string
}

func (c *SpeedBinb) Pages(uri url.URL, imageChan chan<- static.Image, requestOpts *rq.OptsFn) error {
	doc, err := utils.Document(uri.String(), unwrapRequestOpts(requestOpts))
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
		return handleV016113(uri, ptbinb, unwrapRequestOpts(requestOpts), imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") && q.Has("u0") && q.Has("u1") {
		return handleV016452(uri, ptbinb, unwrapRequestOpts(requestOpts), imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") && q.Has("u1") {
		return handleV016201(uri, imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") {
		return handleV016130(uri, ptbinb, unwrapRequestOpts(requestOpts), imageChan)
	} else {
		return handleV016061(uri, imageChan, pagesContent)
	}
}

func unwrapRequestOpts(opts *rq.OptsFn) rq.OptsFn {
	if opts == nil {
		return func(config *rq.Opts) {}
	}

	return *opts
}

func New(domain string) *SpeedBinb {
	return &SpeedBinb{domain: domain}
}
