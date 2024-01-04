package speed_binb

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/pkg/request"
	"net/url"
	"strings"
)

type SpeedBinb struct {
	*connectors.Base
}

func New(domain string) *SpeedBinb {
	return &SpeedBinb{Base: connectors.NewBase(domain)}
}

func (s *SpeedBinb) Context() *connectors.Base {
	return s.Base
}

func (s *SpeedBinb) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	var c = request.Config{Headers: map[string]string{}}
	connectorConfig, exists := config.State.Sites[s.Domain]

	if exists {
		c.Headers["cookie"] = connectorConfig.Session
	}

	doc, err := request.GetDocument(uri.String(), &c)
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
		return handleV016113(uri, imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") && q.Has("u0") && q.Has("u1") {
		return handleV016452(uri, ptbinb, &c, imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") && q.Has("u1") {
		return handleV016201(uri, imageChan)
	} else if ptbinbExists && strings.Contains(ptbinb, "bibGetCntntInfo") {
		return handleV016130(uri, imageChan)
	} else {
		return handleV016061(uri, imageChan, pagesContent)
	}
}
