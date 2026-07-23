package cmoa

import (
	"github.com/sekiju/mdl/extractor/template/speed_binb"
	"github.com/sekiju/mdl/sdk/manga"
	"regexp"
	"resty.dev/v3"
)

type Extractor struct {
	settings *manga.Settings
}

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	return nil, manga.ErrChapterListingUnsupported
}

func (e *Extractor) SupportsChapterListing() bool {
	return false
}

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
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

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	if e.settings.Cookie == nil {
		return nil, manga.ErrCredentialsRequired
	}

	req := resty.New().R().SetHeader("Cookie", *e.settings.Cookie)
	return speed_binb.New(req).FindChapterPages(chapter)
}

func (e *Extractor) SetSettings(settings manga.Settings) {
	e.settings = &settings
}

var re = regexp.MustCompile(`https://www.cmoa.jp/bib/speedreader/[?&]cid=([^&]+)`)

func extractViewerID(URL string) (string, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return "", manga.ErrInvalidURLFormat
	}

	return matches[1], nil
}

func New() (manga.Extractor, error) {
	return &Extractor{settings: &manga.Settings{}}, nil
}
