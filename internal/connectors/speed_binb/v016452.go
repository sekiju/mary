package speed_binb

import (
	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
	"strings"
)

func handleV016452(uri url.URL, apiUrl string, requestOpt request.OptsFn, imageChan chan<- static.Image) error {
	log.Trace().Msg("bimb version 016452")

	cid := uri.Query().Get("cid")
	sharingKey := tt(cid)

	q := uri.Query()
	q.Set("k", sharingKey)
	q.Del("rurl")

	uri.Path = apiUrl
	uri.RawQuery = q.Encode()

	bibGetCntntInfoItems, err := request.Get[BibGetCntntInfo](uri.String(), requestOpt)
	if err != nil {
		return err
	}

	log.Trace().Msgf("bibGetCntntInfo: %s", uri.String())

	if bibGetCntntInfoItems.Body.Result != 1 {
		return fmt.Errorf("invalid bibGetCntntInfoItems result")
	}

	bibGetCntntInfo := bibGetCntntInfoItems.Body.Items[0]

	ctbl := pt(cid, sharingKey, bibGetCntntInfo.Ctbl)
	ptbl := pt(cid, sharingKey, bibGetCntntInfo.Ptbl)
	if bibGetCntntInfo.ServerType != 0 {
		return fmt.Errorf("unsupported speedbinb server type")
	}

	sbcGetCntntUrl, err := url.Parse(utils.JoinURL(bibGetCntntInfo.ContentsServer, "sbcGetCntnt.php"))
	if err != nil {
		return err
	}

	q.Del("k")
	q.Set("p", bibGetCntntInfo.P)
	q.Set("q", "1")
	q.Set("vm", strconv.Itoa(bibGetCntntInfo.ViewMode))
	q.Set("dmytime", bibGetCntntInfo.ContentDate)
	sbcGetCntntUrl.RawQuery = q.Encode()

	log.Trace().Msgf("sbcGetCntntUrl: %s", sbcGetCntntUrl.String())

	sbcGetCntnt, err := request.Get[SbcGetCntnt](sbcGetCntntUrl.String(), nil)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sbcGetCntnt.Body.Ttx))
	if err != nil {
		return err
	}

	sbcGetImgUrl := sbcGetCntntUrl
	sbcGetImgUrl.Path = strings.Replace(sbcGetCntntUrl.Path, "sbcGetCntnt", "sbcGetImg", 1)

	tImages := doc.Find("t-case:first-of-type t-img")
	indexNamer := utils.NewIndexNamer(tImages.Size())

	tImages.Each(func(i int, selection *goquery.Selection) {
		src, _ := selection.Attr("src")

		q = sbcGetImgUrl.Query()
		q.Set("src", src)
		sbcGetImgUrl.RawQuery = q.Encode()

		imgUrl := sbcGetImgUrl.String()
		log.Trace().Msgf("%s: %s", indexNamer.Get(i, ".jpg"), imgUrl)

		process016130(imgUrl, src, indexNamer.Get(i, ".jpg"), ctbl, ptbl, imageChan)
	})

	close(imageChan)
	return nil
}
