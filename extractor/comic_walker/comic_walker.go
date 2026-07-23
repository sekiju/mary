package comic_walker

import (
	"context"
	"encoding/hex"
	"fmt"
	json "github.com/bytedance/sonic"
	"github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/sdk/manga"
	"regexp"
	"resty.dev/v3"
	"strconv"
)

var httpClient = resty.New()

type Extractor struct {
	util.Base
}

type searchFn func(episodeID string) ([]*manga.Chapter, error)

func toChapter(id, code, title string, episodeNo int, workCode string) manga.Chapter {
	return manga.Chapter{
		ID:      id,
		Number:  strconv.Itoa(episodeNo),
		Title:   title,
		Index:   uint(episodeNo - 1),
		URL:     fmt.Sprintf("https://comic-walker.com/detail/%s/episodes/%s", workCode, code),
		MangaID: workCode,
	}
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	parsedURL, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.R().SetContext(ctx).Get(fmt.Sprintf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeType=first", parsedURL.WorkCode))
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, manga.ErrMangaNotFound
	}

	var episodeResult EpisodeResult
	if err = json.Unmarshal(res.Bytes(), &episodeResult); err != nil {
		return nil, err
	}

	firstChapter := toChapter(episodeResult.Episode.Id, episodeResult.Episode.Code, episodeResult.Episode.Title, episodeResult.Episode.Internal.EpisodeNo, parsedURL.WorkCode)
	chapters := []*manga.Chapter{&firstChapter}

	var fn searchFn
	fn = func(episodeID string) ([]*manga.Chapter, error) {
		res, err = httpClient.R().SetContext(ctx).Get(fmt.Sprintf("https://comic-walker.com/api/contents/viewer-jump-forward?episodeId=%s", episodeID))
		if err != nil {
			return nil, err
		}

		var viewerJumpForwardResult ViewerJumpForwardResult
		if err = json.Unmarshal(res.Bytes(), &viewerJumpForwardResult); err != nil {
			return nil, err
		}

		if viewerJumpForwardResult.Episode != nil {
			ep := viewerJumpForwardResult.Episode
			nextChapter := toChapter(ep.Id, ep.Code, ep.Title, ep.Internal.EpisodeNo, parsedURL.WorkCode)
			chapters = append(chapters, &nextChapter)

			return fn(viewerJumpForwardResult.Episode.Id)
		}

		return chapters, nil
	}

	return fn(episodeResult.Episode.Id)
}

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	parsedURL, err := parseURL(URL)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.R().SetContext(ctx).Get(fmt.Sprintf("https://comic-walker.com/api/contents/details/episode?workCode=%s&episodeCode=%s&episodeType=first", parsedURL.WorkCode, parsedURL.EpisodeCode))
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, manga.ErrChapterNotFound
	}

	var episodeResult EpisodeResult
	if err = json.Unmarshal(res.Bytes(), &episodeResult); err != nil {
		return nil, err
	}

	chapter := toChapter(episodeResult.Episode.Id, episodeResult.Episode.Code, episodeResult.Episode.Title, episodeResult.Episode.Internal.EpisodeNo, parsedURL.WorkCode)
	return &chapter, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := httpClient.R().SetContext(ctx).Get(fmt.Sprintf("https://comic-walker.com/api/contents/viewer?episodeId=%s&imageSizeType=width%%3A1284", chapter.ID))
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, manga.ErrChapterNotFound
	}

	var viewerResult ViewerResult
	if err = json.Unmarshal(res.Bytes(), &viewerResult); err != nil {
		return nil, err
	}

	return util.BuildPages(len(viewerResult.Manuscripts), ".webp", func(index int, filename string) (*manga.Page, error) {
		page := viewerResult.Manuscripts[index]
		return &manga.Page{
			Index:    uint(index),
			URL:      page.DrmImageUrl,
			Filename: filename,
			Decode: func(b []byte) ([]byte, error) {
				if len(page.DrmHash) < 16 {
					return nil, fmt.Errorf("comic_walker: drmHash too short: got %d chars, want at least 16: %w", len(page.DrmHash), manga.ErrMalformedChapterData)
				}

				keyBytes, err := hex.DecodeString(page.DrmHash[:16])
				if err != nil {
					return nil, err
				}

				r, i := len(b), len(keyBytes)
				decodedBytes := make([]byte, r)

				for a := 0; a < r; a++ {
					decodedBytes[a] = b[a] ^ keyBytes[a%i]
				}

				return decodedBytes, nil
			},
		}, nil
	})
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: util.Base{Settings: &manga.Settings{}}}, nil
}

var re = regexp.MustCompile("https://comic-walker.com/detail/(KC_[a-zA-Z0-9_]*)(/episodes/(KC_[a-zA-Z0-9_]*))?")

type extractorURL struct {
	WorkCode    string
	EpisodeCode string
}

func parseURL(URL string) (*extractorURL, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) < 2 {
		return nil, manga.ErrInvalidURLFormat
	}

	res := &extractorURL{WorkCode: matches[1]}

	if len(matches) == 4 {
		res.EpisodeCode = matches[3]
	}

	return res, nil
}
