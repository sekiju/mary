package speed_binb

import (
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"image"
	"image/draw"
	"net/url"
	"regexp"
	"strconv"
)

func handleV016061(uri url.URL, imageChan chan<- connectors.ReaderImage, selection *goquery.Selection) error {
	pages := selection.Find("div[data-ptimg$=\"ptimg.json\"]").Map(func(i int, s *goquery.Selection) string {
		text, _ := s.Attr("data-ptimg")
		return text
	})

	if len(pages) > 0 {
		fnf := utils.NewIndexNameFormatter(len(pages))
		for i := 0; i < len(pages); i++ {
			process016061(uri, pages[i], fnf.GetName(i, ".jpg"), imageChan)
		}
	} else {
		return fmt.Errorf("unsupported speedbinb reader")
	}

	close(imageChan)
	return nil
}

func process016061(uri url.URL, page string, fileName string, imageChan chan<- connectors.ReaderImage) {
	var fn connectors.ImageFunction
	fn = func() (image.Image, error) {
		ptimg, err := request.Get[Ptimg](request.JoinURL(uri.String(), page), nil)
		if err != nil {
			return nil, err
		}

		img, err := request.Get[image.Image](request.JoinURL(uri.String(), "data", ptimg.Resources.I.Src), nil)
		if err != nil {
			return nil, err
		}

		return descramble016061(img, ptimg.Views), nil
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}

func descramble016061(img image.Image, views []PtimgView) image.Image {
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

	return descrambledImg
}
