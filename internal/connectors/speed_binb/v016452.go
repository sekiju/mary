package speed_binb

import (
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"image"
	"image/draw"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func handleV016452(uri url.URL, apiUrl string, c *request.Config, imageChan chan<- connectors.ReaderImage) error {
	cid := uri.Query().Get("cid")
	sharingKey := tt(cid)

	q := uri.Query()
	q.Set("k", sharingKey)
	q.Del("rurl")

	uri.Path = apiUrl
	uri.RawQuery = q.Encode()

	bibGetCntntInfoItems, err := request.Get[BibGetCntntInfo](uri.String(), c)
	if err != nil {
		return err
	}

	if bibGetCntntInfoItems.Result != 1 {
		return fmt.Errorf("invalid bibGetCntntInfoItems result")
	}

	bibGetCntntInfo := bibGetCntntInfoItems.Items[0]

	ctbl := pt(cid, sharingKey, bibGetCntntInfo.Ctbl)
	ptbl := pt(cid, sharingKey, bibGetCntntInfo.Ptbl)
	if bibGetCntntInfo.ServerType != 0 {
		return fmt.Errorf("unsupported speedbinb server type")
	}

	sbcGetCntntUrl, err := url.Parse(request.JoinURL(bibGetCntntInfo.ContentsServer, "sbcGetCntnt.php"))
	if err != nil {
		return err
	}

	q.Del("k")
	q.Set("p", bibGetCntntInfo.P)
	q.Set("q", "1")
	q.Set("vm", strconv.Itoa(bibGetCntntInfo.ViewMode))
	q.Set("dmytime", bibGetCntntInfo.ContentDate)
	sbcGetCntntUrl.RawQuery = q.Encode()

	sbcGetCntnt, err := request.Get[SbcGetCntnt](sbcGetCntntUrl.String(), nil)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(sbcGetCntnt.Ttx))
	if err != nil {
		return err
	}

	sbcGetImgUrl := sbcGetCntntUrl
	sbcGetImgUrl.Path = strings.Replace(sbcGetCntntUrl.Path, "sbcGetCntnt", "sbcGetImg", 1)

	tImages := doc.Find("t-case:first-of-type t-img")
	fileName := utils.NewIndexNameFormatter(tImages.Size())

	tImages.Each(func(i int, selection *goquery.Selection) {
		src, _ := selection.Attr("src")
		process016452(*sbcGetImgUrl, fileName.GetName(i, ".jpg"), src, ctbl, ptbl, imageChan)
	})

	close(imageChan)
	return nil
}

// todo: move to 0161130
func process016452(uri url.URL, fileName, src string, ctbl, ptbl []string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		q := uri.Query()
		q.Set("src", src)
		uri.RawQuery = q.Encode()

		prototype := lt(src, ctbl, ptbl)
		if prototype == nil || !prototype.vt() {
			return nil, fmt.Errorf("prototype.vt() dont exists")
		}

		img, err := request.Get[image.Image](uri.String(), nil)
		if err != nil {
			return nil, err
		}

		e := prototype.dt(img.Bounds())

		view := DescrambleView{Width: e.Dx(), Height: e.Dy(), Transfers: []DescrambleTransfer{{0, prototype.gt(img.Bounds())}}}

		return descramble016452(img, &view), nil
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}

// todo: move to 0161130
func descramble016452(img image.Image, view *DescrambleView) image.Image {
	descrambledImg := image.NewRGBA(image.Rect(0, 0, view.Width, view.Height))

	for _, part := range view.Transfers[0].Coords {
		wherePlaceRect := image.Rect(part.XDest, part.YDest, part.XDest+part.Width, part.YDest+part.Height)
		whereTakeRect := image.Rect(part.XSrc, part.YSrc, part.XSrc+part.Width, part.YSrc+part.Height)

		draw.Draw(descrambledImg, wherePlaceRect, img, whereTakeRect.Min, draw.Src)
	}

	return descrambledImg
}

// todo: move to 0161130
// todo: rename tt to resolveKey/generateKey

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
