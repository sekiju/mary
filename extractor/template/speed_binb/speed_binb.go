package speed_binb

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	json "github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/sdk/manga"
	"github.com/sekiju/mdl/sdk/manga/pluginutil"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"net/url"
	"regexp"
	"resty.dev/v3"
	"strconv"
	"strings"
)

type Extractor struct {
	req *resty.Request
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := e.req.SetContext(ctx).Get(chapter.URL)
	if err != nil || res.StatusCode() != 200 {
		return nil, manga.ErrChapterNotFound
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Bytes()))
	if err != nil {
		return nil, err
	}

	content := doc.Find("div#content.pages").First()
	apiURL, apiURLExists := content.Attr("data-ptbinb")

	parsedURL, err := url.Parse(chapter.URL)
	if err != nil {
		return nil, err
	}

	query := parsedURL.Query()

	switch {
	case apiURLExists && strings.Contains(apiURL, "bibGetCntntInfo") && query.Has("u1"):
		return e.v016452(ctx, parsedURL, apiURL)
	default:
		return e.v016061(ctx, parsedURL, content)
	}
}

func (e *Extractor) v016061(ctx context.Context, parsedURL *url.URL, content *goquery.Selection) ([]*manga.Page, error) {
	log.Trace().Msg("SpeedBinb version: v016061")

	tPages := content.Find("div[data-ptimg$=\"ptimg.json\"]").Map(func(i int, s *goquery.Selection) string {
		text, _ := s.Attr("data-ptimg")
		return text
	})

	if len(tPages) > 0 {
		return pluginutil.BuildPages(len(tPages), ".png", func(i int, filename string) (*manga.Page, error) {
			src := tPages[i]
			res, err := resty.New().R().SetContext(ctx).Get(parsedURL.String() + "/" + src)
			if err != nil {
				return nil, err
			}

			var ptImg Ptimg
			if err = json.Unmarshal(res.Bytes(), &ptImg); err != nil {
				return nil, err
			}

			return &manga.Page{
				URL:      parsedURL.String() + "/data/" + ptImg.Resources.I.Src,
				Filename: filename,
				Index:    uint(i),
				Decode:   e.decode016061(ptImg.Views),
			}, nil
		})
	} else {
		return nil, manga.ErrMethodUnimplemented
	}
}

var reDecode016061 = regexp.MustCompile("[:,+>]")

// descrambleImage rebuilds a width x height image from img by copying each
// coords entry from its source rect into its destination rect, then
// re-encodes the result as PNG. Shared by decode016061 and decode016130,
// which differ only in how they compute width/height/coords.
func descrambleImage(img image.Image, width, height int, coords []DescrambleCord) ([]byte, error) {
	descrambledImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for _, part := range coords {
		dstRect := image.Rect(part.XDest, part.YDest, part.XDest+part.Width, part.YDest+part.Height)
		srcRect := image.Rect(part.XSrc, part.YSrc, part.XSrc+part.Width, part.YSrc+part.Height)
		draw.Draw(descrambledImg, dstRect, img, srcRect.Min, draw.Src)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, descrambledImg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *Extractor) decode016061(views []PtimgView) manga.DecodeFunc {
	coords := make([]DescrambleCord, 0, len(views[0].Coords))
	for _, part := range views[0].Coords {
		num := reDecode016061.Split(part, -1)

		sourceX, _ := strconv.Atoi(num[1])
		sourceY, _ := strconv.Atoi(num[2])
		partWidth, _ := strconv.Atoi(num[3])
		partHeight, _ := strconv.Atoi(num[4])
		targetX, _ := strconv.Atoi(num[5])
		targetY, _ := strconv.Atoi(num[6])

		coords = append(coords, DescrambleCord{
			XSrc: sourceX, YSrc: sourceY,
			Width: partWidth, Height: partHeight,
			XDest: targetX, YDest: targetY,
		})
	}

	return func(b []byte) ([]byte, error) {
		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}

		return descrambleImage(img, views[0].Width, views[0].Height, coords)
	}
}

func (e *Extractor) v016452(ctx context.Context, parsedURL *url.URL, apiURL string) ([]*manga.Page, error) {
	log.Trace().Msg("SpeedBinb version: v016452")

	cid := parsedURL.Query().Get("cid")
	sharingKey := generateSharingKey(cid)

	query := parsedURL.Query()
	query.Set("k", sharingKey)
	query.Del("rurl")

	parsedURL.Path = apiURL
	parsedURL.RawQuery = query.Encode()

	res, err := e.req.SetContext(ctx).Get(parsedURL.String())
	if err != nil {
		return nil, err
	}

	var bibGetCntntInfoItems BibGetCntntInfo
	if err = json.Unmarshal(res.Bytes(), &bibGetCntntInfoItems); err != nil {
		return nil, err
	}

	if bibGetCntntInfoItems.Result != 1 {
		return nil, fmt.Errorf("invalid bibGetCntntInfoItems result: %w", manga.ErrMalformedChapterData)
	}

	bibGetCntntInfo := bibGetCntntInfoItems.Items[0]

	if bibGetCntntInfo.ServerType != 0 {
		return nil, fmt.Errorf("unsupported speedbinb server type %d: %w", bibGetCntntInfo.ServerType, manga.ErrUnsupportedContentFormat)
	}

	sbcGetCntntUrl, err := url.Parse(bibGetCntntInfo.ContentsServer + "/sbcGetCntnt.php")
	if err != nil {
		return nil, err
	}

	query.Del("k")
	query.Set("p", bibGetCntntInfo.P)
	query.Set("q", "1")
	query.Set("vm", strconv.Itoa(bibGetCntntInfo.ViewMode))
	query.Set("dmytime", bibGetCntntInfo.ContentDate)
	sbcGetCntntUrl.RawQuery = query.Encode()

	res, err = e.req.SetContext(ctx).Get(sbcGetCntntUrl.String())
	if err != nil {
		return nil, err
	}

	var sbcGetCntn SbcGetCntnt
	if err = json.Unmarshal(res.Bytes(), &sbcGetCntn); err != nil {
		return nil, err
	}

	tDoc, err := goquery.NewDocumentFromReader(strings.NewReader(sbcGetCntn.Ttx))
	if err != nil {
		return nil, err
	}

	ctbl := decodeTable(cid, sharingKey, bibGetCntntInfo.Ctbl)
	ptbl := decodeTable(cid, sharingKey, bibGetCntntInfo.Ptbl)

	sbcGetImgUrl := sbcGetCntntUrl
	sbcGetImgUrl.Path = strings.Replace(sbcGetCntntUrl.Path, "sbcGetCntnt", "sbcGetImg", 1)

	tImages := tDoc.Find("t-case:first-of-type t-img")

	return pluginutil.BuildPages(tImages.Length(), ".png", func(i int, filename string) (*manga.Page, error) {
		src, _ := tImages.Eq(i).Attr("src")

		query = sbcGetImgUrl.Query()
		query.Set("src", src)
		sbcGetImgUrl.RawQuery = query.Encode()

		return &manga.Page{
			URL:      sbcGetImgUrl.String(),
			Filename: filename,
			Index:    uint(i),
			Decode:   e.decode016130(src, ctbl, ptbl),
		}, nil
	})
}

func (e *Extractor) decode016130(imgSrc string, ctbl, ptbl []string) manga.DecodeFunc {
	return func(b []byte) ([]byte, error) {
		prototype := selectPrototype(imgSrc, ctbl, ptbl)
		if prototype == nil || !prototype.vt() {
			return nil, fmt.Errorf("prototype.vt() dont exists: %w", manga.ErrUnsupportedContentFormat)
		}

		img, _, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			return nil, err
		}

		destRect := prototype.dt(img.Bounds())
		coords := prototype.gt(img.Bounds())

		return descrambleImage(img, destRect.Dx(), destRect.Dy(), coords)
	}
}

func New(requests ...*resty.Request) *Extractor {
	var req *resty.Request
	if len(requests) > 0 {
		req = requests[0]
	} else {
		req = resty.New().R()
	}

	return &Extractor{req}
}
