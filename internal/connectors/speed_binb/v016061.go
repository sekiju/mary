package speed_binb

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/url"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"

	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
)

func handleV016061(uri url.URL, imageChan chan<- static.Image, selection *goquery.Selection) error {
	log.Trace().Msg("bimb version 016061")

	pages := selection.Find("div[data-ptimg$=\"ptimg.json\"]").Map(func(i int, s *goquery.Selection) string {
		text, _ := s.Attr("data-ptimg")
		return text
	})

	if len(pages) > 0 {
		indexNamer := utils.NewIndexNamer(len(pages))
		for i := 0; i < len(pages); i++ {
			var fn static.ImageFn
			fn = func() ([]byte, error) {
				ptimgResponse, err := request.Get[Ptimg](utils.JoinURL(uri.String(), pages[i]))
				if err != nil {
					return nil, err
				}

				imageResponse, err := request.Get[image.Image](utils.JoinURL(uri.String(), "data", ptimgResponse.Body.Resources.I.Src))
				if err != nil {
					return nil, err
				}

				return descramble016061(imageResponse.Body, ptimgResponse.Body.Views), nil
			}

			imageChan <- static.NewImage(indexNamer.Get(i, ".jpg"), &fn)
		}
	} else {
		return fmt.Errorf("unsupported speedbinb reader")
	}

	close(imageChan)
	return nil
}

func descramble016061(img image.Image, views []PtimgView) []byte {
	descrambledImg := image.NewRGBA(image.Rect(0, 0, views[0].Width, views[0].Height))

	re := regexp.MustCompile("[:,+>]")
	for _, part := range views[0].Coords {
		num := re.Split(part, -1)

		sourceX, _ := strconv.Atoi(num[1])
		sourceY, _ := strconv.Atoi(num[2])
		partWidth, _ := strconv.Atoi(num[3])
		partHeight, _ := strconv.Atoi(num[4])
		targetX, _ := strconv.Atoi(num[5])
		targetY, _ := strconv.Atoi(num[6])

		dstRect, drawRect := image.Rect(targetX, targetY, targetX+partWidth, targetY+partHeight), image.Rect(sourceX, sourceY, sourceX+partWidth, sourceY+partHeight)
		draw.Draw(descrambledImg, dstRect, img, drawRect.Min, draw.Src)
	}

	buf := new(bytes.Buffer)
	_ = jpeg.Encode(buf, descrambledImg, nil)

	return buf.Bytes()
}
