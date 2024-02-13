package speed_binb

import (
	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"image"
	"image/draw"
	"image/jpeg"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func handleV016130(uri url.URL, apiUrl string, imageChan chan<- static.Image) error {
	log.Trace().Msg("bimb version 016130")

	cid := uri.Query().Get("cid")
	sharingKey := tt(cid)

	path, q := utils.ExportURLQueries(apiUrl)
	uri.RawQuery = ""
	uri.Path = path

	q.Set("k", sharingKey)
	q.Set("cid", cid)
	uri.RawQuery = q.Encode()

	log.Trace().Msgf("bibGetCntntInfo: %s", uri.String())

	bibGetCntntInfoItems, err := request.Get[BibGetCntntInfo16130](uri.String())
	if err != nil {
		return err
	}

	if bibGetCntntInfoItems.Body.Result != 1 {
		return fmt.Errorf("invalid bibGetCntntInfoItems result")
	}

	bibGetCntntInfo := bibGetCntntInfoItems.Body.Items[0]

	ctbl := pt(cid, sharingKey, bibGetCntntInfo.Ctbl)
	ptbl := pt(cid, sharingKey, bibGetCntntInfo.Ptbl)
	if bibGetCntntInfo.ServerType == 0 {
		return fmt.Errorf("unsupported speedbinb server type")
	}

	if bibGetCntntInfo.ServerType == 1 {
		return fmt.Errorf("unsupported speedbinb server type")
	}

	if bibGetCntntInfo.ServerType == 2 {
		sbcGetCntntUrl, err := url.Parse(utils.JoinURL(bibGetCntntInfo.ContentsServer, "content"))
		if err != nil {
			return err
		}

		//q.Set("dmytime", bibGetCntntInfo.ContentDate)
		//sbcGetCntntUrl.RawQuery = q.Encode()

		log.Trace().Msgf("sbcGetCntntUrl: %s", sbcGetCntntUrl.String())

		sbcGetCntnt, err := request.Get[SbcGetCntnt](sbcGetCntntUrl.String())
		if err != nil {
			return err
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(sbcGetCntnt.Body.Ttx))
		if err != nil {
			return err
		}

		sbcGetImgUrl := sbcGetCntntUrl
		sbcGetImgUrl.Path = strings.Replace(sbcGetCntntUrl.Path, "content", "img", 1)

		tImages := doc.Find("t-case:first-of-type t-img")
		indexNamer := utils.NewIndexNamer(tImages.Size())

		tImages.Each(func(i int, selection *goquery.Selection) {
			src, _ := selection.Attr("src")

			imgUrl := utils.JoinURL(sbcGetImgUrl.String(), src)
			log.Trace().Msgf("%s: %s", indexNamer.Get(i, ".jpg"), imgUrl)

			process016130(imgUrl, src, indexNamer.Get(i, ".jpg"), ctbl, ptbl, imageChan)
		})
	}

	close(imageChan)
	return nil
}

func process016130(imgUrl, imgSrc, fileName string, ctbl, ptbl []string, imageChan chan<- static.Image) {
	var fn static.ImageFn
	fn = func() ([]byte, error) {
		prototype := lt(imgSrc, ctbl, ptbl)
		if prototype == nil || !prototype.vt() {
			return nil, fmt.Errorf("prototype.vt() dont exists")
		}

		img, err := request.Get[image.Image](imgUrl)
		if err != nil {
			return nil, err
		}

		e := prototype.dt(img.Body.Bounds())

		view := DescrambleView{Width: e.Dx(), Height: e.Dy(), Transfers: []DescrambleTransfer{{0, prototype.gt(img.Body.Bounds())}}}

		return descramble016130(img.Body, &view), nil
	}

	imageChan <- static.NewImage(fileName, &fn)
}

func descramble016130(img image.Image, view *DescrambleView) []byte {
	descrambledImg := image.NewRGBA(image.Rect(0, 0, view.Width, view.Height))

	for _, part := range view.Transfers[0].Coords {
		wherePlaceRect := image.Rect(part.XDest, part.YDest, part.XDest+part.Width, part.YDest+part.Height)
		whereTakeRect := image.Rect(part.XSrc, part.YSrc, part.XSrc+part.Width, part.YSrc+part.Height)

		draw.Draw(descrambledImg, wherePlaceRect, img, whereTakeRect.Min, draw.Src)
	}

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, descrambledImg, nil)

	return buf.Bytes()
}

func padStart(input string, length int, padChar string) string {
	if len(input) >= length {
		return input
	}

	padding := strings.Repeat(padChar, length-len(input))
	return padding + input
}

func tt(t string) string {
	timestampMilliseconds := time.Now().UnixNano() / int64(time.Millisecond)
	n := strconv.FormatInt(timestampMilliseconds, 16)
	nHex := padStart(n, 16, "x")
	repeatedString := strings.Repeat(t, (16+len(t)-1)/len(t))
	r := repeatedString[:16]
	e := repeatedString[len(repeatedString)-16:]

	s, u, h := 0, 0, 0
	result := ""
	for i, char := range nHex {
		s ^= int(char)
		u ^= int(r[i])
		h ^= int(e[i])

		result += string(char)
		result += string("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"[(s+u+h)&63])
	}

	return result
}

func pt(t, i, n string) []string {
	r := t + ":" + i
	e := 0
	for s := 0; s < len(r); s++ {
		e += int(r[s]) << (s % 16)
	}
	e &= 2147483647
	if e == 0 {
		e = 305419896
	}

	u := ""
	h := e
	for s := 0; s < len(n); s++ {
		h = h>>1 ^ 1210056708&-(1&h)
		o := (int(n[s])-32+h)%94 + 32
		u += string(rune(o))
	}

	var result []string
	err := json.Unmarshal([]byte(u), &result)
	if err != nil {
		return nil
	}
	return result
}

func lastIndexOf(s string, char rune) int {
	for i := len(s) - 1; i >= 0; i-- {
		if rune(s[i]) == char {
			return i
		}
	}
	return -1
}

func lt(t string, ctbl, ptbl []string) SpeedBinbDecoder {
	i := [2]int{0, 0}

	if t != "" {
		n := lastIndexOf(t, '/') + 1
		r := len(t) - n
		for e := 0; e < r; e++ {
			i[e%2] += int(t[e+n])
		}

		i[0] %= 8
		i[1] %= 8
	}

	s, u := ptbl[i[0]], ctbl[i[1]]

	if s[0] == '=' && u[0] == '=' {
		return NewSpeedbinbF(u, s)
	} else if regexp.MustCompile("^[0-9]").MatchString(u) && regexp.MustCompile("^[0-9]").MatchString(s) {
		return NewSpeedbinbA(u, s)
	} else if "" == u && "" == s {
		return NewSpeedbinbH()
	}

	return nil
}
