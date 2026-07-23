package storia_takeshobo

import (
	"bytes"
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/sekiju/mdl/extractor/template/speed_binb"
	"github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/sdk/manga"
	"regexp"
	"resty.dev/v3"
)

type Extractor struct {
	util.Base
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	return nil, manga.ErrChapterListingUnsupported
}

func (e *Extractor) SupportsChapterListing() bool {
	return false
}

var re = regexp.MustCompile(`https://storia.takeshobo.co.jp/_files/([a-zA-Z0-9_]*)/(\d*)`)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	res, err := resty.New().R().SetContext(ctx).Get(URL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Bytes()))
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      matches[2],
		Number:  "",
		Title:   doc.Find("title").Text(),
		Index:   0,
		URL:     URL,
		MangaID: matches[1],
	}, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	return speed_binb.New().FindChapterPages(ctx, chapter)
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: util.Base{Settings: &manga.Settings{}}}, nil
}
