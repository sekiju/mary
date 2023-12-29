package takeshobo

import (
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"github.com/PuerkitoBio/goquery"
	"image"
	"image/draw"
	"net/url"
	"regexp"
	"strconv"
)

type StoriaTakeshobo struct {
	*connectors.Base
}

func New() *StoriaTakeshobo {
	return &StoriaTakeshobo{Base: connectors.NewBase("storia.takeshobo.co.jp")}
}

func (t *StoriaTakeshobo) Context() *connectors.Base {
	return t.Base
}

func (t *StoriaTakeshobo) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	doc, err := request.GetDocument(uri.String(), nil)
	if err != nil {
		return err
	}

	var pages []string
	doc.Find("div#content.pages > div").Each(func(i int, s *goquery.Selection) {
		text, _ := s.Attr("data-ptimg")
		pages = append(pages, text)
	})

	fnf := utils.NewIndexNameFormatter(len(pages))
	for i := 0; i < len(pages); i++ {
		processPage(uri, pages[i], fnf.GetName(i, ".jpg"), imageChan)
	}

	close(imageChan)

	return nil
}

func processPage(uri url.URL, page string, fileName string, imageChan chan<- connectors.ReaderImage) {
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

		return descrambleImage(img, ptimg.Views), nil
	}

	imageChan <- connectors.NewConnectorImage(fileName, &fn)
}

func descrambleImage(img image.Image, views []PtimgView) image.Image {
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
