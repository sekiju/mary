package speed_binb

import (
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"image"
	"image/draw"
	"net/url"
	"strings"
)

func handleV016130(uri url.URL, apiUrl string, imageChan chan<- connectors.ReaderImage) error {
	cid := uri.Query().Get("cid")
	sharingKey := tt(cid)

	path, q := request.ExportURLQueries(apiUrl)
	uri.RawQuery = ""
	uri.Path = path

	q.Set("k", sharingKey)
	q.Set("cid", cid)
	uri.RawQuery = q.Encode()

	log.Trace().Msgf("bibGetCntntInfo: %s", uri.String())

	bibGetCntntInfoItems, err := request.Get[BibGetCntntInfo16130](uri.String(), nil)
	if err != nil {
		return err
	}

	if bibGetCntntInfoItems.Result != 1 {
		return fmt.Errorf("invalid bibGetCntntInfoItems result")
	}

	bibGetCntntInfo := bibGetCntntInfoItems.Items[0]

	ctbl := pt(cid, sharingKey, bibGetCntntInfo.Ctbl)
	ptbl := pt(cid, sharingKey, bibGetCntntInfo.Ptbl)
	if bibGetCntntInfo.ServerType == 0 {
		return fmt.Errorf("unsupported speedbinb server type")
	}

	if bibGetCntntInfo.ServerType == 1 {
		return fmt.Errorf("unsupported speedbinb server type")
	}

	if bibGetCntntInfo.ServerType == 2 {
		sbcGetCntntUrl, err := url.Parse(request.JoinURL(bibGetCntntInfo.ContentsServer, "content"))
		if err != nil {
			return err
		}

		//q.Set("dmytime", bibGetCntntInfo.ContentDate)
		//sbcGetCntntUrl.RawQuery = q.Encode()

		log.Trace().Msgf("sbcGetCntntUrl: %s", sbcGetCntntUrl.String())

		sbcGetCntnt, err := request.Get[SbcGetCntnt](sbcGetCntntUrl.String(), nil)
		if err != nil {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(sbcGetCntnt.Ttx))
		if err != nil {
			return err
		}

		sbcGetImgUrl := sbcGetCntntUrl
		sbcGetImgUrl.Path = strings.Replace(sbcGetCntntUrl.Path, "content", "img", 1)

		tImages := doc.Find("t-case:first-of-type t-img")
		fileName := utils.NewIndexNameFormatter(tImages.Size())

		tImages.Each(func(i int, selection *goquery.Selection) {
			src, _ := selection.Attr("src")

			imgUrl := request.JoinURL(sbcGetImgUrl.String(), src)
			log.Trace().Msgf("%s: %s", fileName.GetName(i, ".jpg"), imgUrl)

			process016130(imgUrl, src, fileName.GetName(i, ".jpg"), ctbl, ptbl, imageChan)
		})
	}

	close(imageChan)
	return nil
}

func process016130(imgUrl, imgSrc, fileName string, ctbl, ptbl []string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		prototype := lt(imgSrc, ctbl, ptbl)
		if prototype == nil || !prototype.vt() {
			return nil, fmt.Errorf("prototype.vt() dont exists")
		}

		img, err := request.Get[image.Image](imgUrl, nil)
		if err != nil {
			return nil, err
		}

		e := prototype.dt(img.Bounds())

		view := DescrambleView{Width: e.Dx(), Height: e.Dy(), Transfers: []DescrambleTransfer{{0, prototype.gt(img.Bounds())}}}

		return descramble016130(img, &view), nil
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}

func descramble016130(img image.Image, view *DescrambleView) image.Image {
	descrambledImg := image.NewRGBA(image.Rect(0, 0, view.Width, view.Height))

	for _, part := range view.Transfers[0].Coords {
		wherePlaceRect := image.Rect(part.XDest, part.YDest, part.XDest+part.Width, part.YDest+part.Height)
		whereTakeRect := image.Rect(part.XSrc, part.YSrc, part.XSrc+part.Width, part.YSrc+part.Height)

		draw.Draw(descrambledImg, wherePlaceRect, img, whereTakeRect.Min, draw.Src)
	}

	return descrambledImg
}
