package giga_viewer

import (
	"github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/sdk/manga"
	"resty.dev/v3"
	"strconv"
)

var httpClient = resty.New()

type Extractor struct {
	settings *manga.Settings
}

type searchFn func(URL string) ([]*manga.Chapter, error)

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	chapters := make([]*manga.Chapter, 0)
	visitedIDs := make(map[string]bool)

	var fn searchFn
	fn = func(episodeURL string) ([]*manga.Chapter, error) {
		if visitedIDs[episodeURL] {
			return nil, nil
		}

		visitedIDs[episodeURL] = true

		res, err := httpClient.R().Get(episodeURL)
		if err != nil {
			return nil, err
		}

		if res.StatusCode() == 404 {
			return nil, manga.ErrMangaNotFound
		}

		html := res.String()

		episode, err := util.ExtractJSONFromHTML[episodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
		if err != nil {
			return nil, err
		}

		chapters = append(chapters, &manga.Chapter{
			ID:      episode.ReadableProduct.Id,
			Number:  strconv.Itoa(episode.ReadableProduct.Number),
			Title:   episode.ReadableProduct.Title,
			Index:   uint(episode.ReadableProduct.Number - 1),
			URL:     episode.ReadableProduct.Permalink,
			MangaID: episode.ReadableProduct.Id,
		})

		if prevURI := episode.ReadableProduct.PrevReadableProductUri; prevURI != nil {
			_, err = fn(*prevURI)
			if err != nil {
				return nil, err
			}
		}

		if nextURI := episode.ReadableProduct.NextReadableProductUri; nextURI != nil {
			_, err = fn(*nextURI)
			if err != nil {
				return nil, err
			}
		}

		return chapters, nil
	}

	return fn(URL)
}

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	res, err := httpClient.R().Get(URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, manga.ErrChapterNotFound
	}

	html := res.String()

	episode, err := util.ExtractJSONFromHTML[episodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
	if err != nil {
		return nil, err
	}

	return &manga.Chapter{
		ID:      episode.ReadableProduct.Id,
		Number:  strconv.Itoa(episode.ReadableProduct.Number),
		Title:   episode.ReadableProduct.Title,
		Index:   uint(episode.ReadableProduct.Number - 1),
		URL:     episode.ReadableProduct.Permalink,
		MangaID: episode.ReadableProduct.Id,
	}, nil
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	req := httpClient.R().SetHeader("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Mobile/15E148 Safari/604.1")
	if e.settings.Cookie != nil {
		req.SetHeader("Cookie", *e.settings.Cookie)
	}

	res, err := req.Get(chapter.URL)

	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, manga.ErrChapterNotFound
	}

	html := res.String()

	episode, err := util.ExtractJSONFromHTML[episodeResult](html, `<script id='episode-json' type='text/json' data-value='`, `'></script>`)
	if err != nil {
		return nil, err
	}

	if !episode.ReadableProduct.IsPublic && !episode.ReadableProduct.HasPurchased {
		return nil, manga.ErrPaidChapter
	}

	var mainPages []*episodeResultPage
	for _, page := range episode.ReadableProduct.PageStructure.Pages {
		if page.Type != "main" {
			continue
		}

		mainPages = append(mainPages, &page)
	}

	chapterPages := make([]*manga.Page, len(mainPages))
	padRenamer := renamer.New(len(mainPages))

	for index, page := range mainPages {
		chapterPages[index] = &manga.Page{
			Index:    uint(index),
			URL:      page.Src,
			Filename: padRenamer.Name(index, ".jpg"),
		}
	}

	return chapterPages, nil
}

func (e *Extractor) SetSettings(settings manga.Settings) {
	e.settings = &settings
}

func New() (manga.Extractor, error) {
	return &Extractor{settings: &manga.Settings{}}, nil
}
