package cmoa

import (
	"context"

	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/extractor/template/speed_binb"
	"github.com/knst0/mdl/sdk/manga"
	"github.com/knst0/mdl/sdk/manga/pluginutil"
	"regexp"

	"resty.dev/v3"
)

type Extractor struct {
	pluginutil.Base
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	return nil, manga.ErrChapterListingUnsupported
}

func (e *Extractor) SupportsChapterListing() bool {
	return false
}

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	ID, err := extractViewerID(URL)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      ID,
		Number:  "",
		Title:   "",
		Index:   0,
		URL:     URL,
		MangaID: "",
	}, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	if e.Settings.Cookie == nil {
		return nil, manga.ErrCredentialsRequired
	}

	req := resty.New().R().SetContext(ctx)
	e.ApplyCookie(req)
	return speed_binb.New(req).FindChapterPages(ctx, chapter)
}

var re = regexp.MustCompile(`https://www.cmoa.jp/bib/speedreader/[?&]cid=([^&]+)`)

func extractViewerID(URL string) (string, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return "", manga.ErrInvalidURLFormat
	}

	return matches[1], nil
}

func init() {
	registry.Register("www.cmoa.jp", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
