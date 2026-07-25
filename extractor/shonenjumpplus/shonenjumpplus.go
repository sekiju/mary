package shonenjumpplus

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"image"
	"image/draw"
	_ "image/jpeg" // NOTE: register JPEG decoder
	"image/png"
	"regexp"
	"sort"
	"strconv"
	"strings"

	json "github.com/bytedance/sonic"
	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/sdk/manga"
	"github.com/knst0/mdl/sdk/manga/pluginutil"

	"resty.dev/v3"
)

// iPhone UA avoids the puzzle-piece scrambling served to desktop browsers.
const iPhoneUA = "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"

var httpClient = resty.New()

type Extractor struct {
	pluginutil.Base
}

var episodeRe = regexp.MustCompile(`https://shonenjumpplus\.com/episode/(\d+)`)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := episodeRe.FindStringSubmatch(URL)
	if len(matches) != 2 {
		return nil, manga.ErrInvalidChapterURL
	}
	episodeID := matches[1]

	ej, err := e.fetchEpisodeJSON(ctx, URL)
	if err != nil {
		return nil, err
	}

	rp := ej.ReadableProduct
	return &manga.Chapter{
		ID:      episodeID,
		Number:  strconv.Itoa(rp.Number),
		Title:   rp.Title,
		URL:     URL,
		MangaID: rp.Series.ID,
	}, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	ej, err := e.fetchEpisodeJSON(ctx, chapter.URL)
	if err != nil {
		return nil, err
	}

	all := ej.ReadableProduct.PageStructure.Pages
	pages := make([]pageEntry, 0, len(all))
	for _, p := range all {
		if p.Src != "" {
			pages = append(pages, p)
		}
	}
	if len(pages) == 0 {
		return nil, fmt.Errorf("%w: no pages found for chapter %s", manga.ErrMalformedChapterData, chapter.ID)
	}

	return pluginutil.BuildPages(len(pages), ".png", func(index int, filename string) (*manga.Page, error) {
		return &manga.Page{
			Index:    uint(index),
			URL:      pages[index].Src,
			Filename: filename,
			Headers:  map[string]string{"User-Agent": iPhoneUA},
			Decode:   descramblePuzzleImage,
		}, nil
	})
}

// puzzleDivideNum and puzzleMultiple mirror the GigaViewer webpack bundle's
// puzzle presenter (DIVIDE_NUM=4, MULTIPLE=8): each page image is split into
// a 4x4 grid of 8px-aligned cells, then cells (row, col) and (col, row) are
// swapped (a grid transpose) before being served to non-browser clients.
const (
	puzzleDivideNum = 4
	puzzleMultiple  = 8
)

// descramblePuzzleImage undoes the GigaViewer puzzle scramble by transposing
// the 4x4 cell grid back into place, then re-encodes as PNG.
func descramblePuzzleImage(raw []byte) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("shonenjumpplus: decode image: %w", err)
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	cellW := (w / (puzzleDivideNum * puzzleMultiple)) * puzzleMultiple
	cellH := (h / (puzzleDivideNum * puzzleMultiple)) * puzzleMultiple

	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	for e := range puzzleDivideNum * puzzleDivideNum {
		row, col := e/puzzleDivideNum, e%puzzleDivideNum
		dstX, dstY := col*cellW, row*cellH

		src := col*puzzleDivideNum + row // grid transpose
		srcCol, srcRow := src%puzzleDivideNum, src/puzzleDivideNum
		srcX, srcY := srcCol*cellW, srcRow*cellH

		for y := range cellH {
			for x := range cellW {
				dst.Set(dstX+x, dstY+y, img.At(bounds.Min.X+srcX+x, bounds.Min.Y+srcY+y))
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, fmt.Errorf("shonenjumpplus: encode png: %w", err)
	}
	return buf.Bytes(), nil
}

var seriesRe = regexp.MustCompile(`https://shonenjumpplus\.com/titles/(\d+)`)

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	var aggregateID string

	if matches := seriesRe.FindStringSubmatch(URL); len(matches) == 2 {
		aggregateID = matches[1]
	} else if matches := episodeRe.FindStringSubmatch(URL); len(matches) == 2 {
		ej, err := e.fetchEpisodeJSON(ctx, URL)
		if err != nil {
			return nil, err
		}
		aggregateID = ej.ReadableProduct.Series.ID
	} else {
		return nil, manga.ErrInvalidURLFormat
	}

	if aggregateID == "" {
		return nil, manga.ErrMangaNotFound
	}

	entries, err := e.fetchAllPaginationEntries(ctx, aggregateID)
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, manga.ErrChapterNotFound
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ReadableProductID < entries[j].ReadableProductID
	})

	chapters := make([]*manga.Chapter, 0, len(entries))
	for i, entry := range entries {
		chapters = append(chapters, &manga.Chapter{
			ID:      entry.ReadableProductID,
			Number:  strconv.Itoa(i + 1),
			Title:   entry.Title,
			Index:   uint(i),
			URL:     entry.ViewerURI,
			MangaID: aggregateID,
		})
	}

	return chapters, nil
}

const paginationPageSize = 50

func (e *Extractor) fetchAllPaginationEntries(ctx context.Context, aggregateID string) ([]paginationEntry, error) {
	var entries []paginationEntry
	for offset := 0; ; offset += paginationPageSize {
		req := httpClient.R().SetContext(ctx).
			SetHeader("Referer", "https://shonenjumpplus.com/")
		e.ApplyCookie(req)

		res, err := req.SetQueryParams(map[string]string{
			"type":         "episode",
			"aggregate_id": aggregateID,
			"offset":       strconv.Itoa(offset),
			"limit":        strconv.Itoa(paginationPageSize),
			"sort_order":   "asc",
			"is_guest":     "1",
		}).Get("https://shonenjumpplus.com/api/viewer/pagination_readable_products")
		if err != nil {
			return nil, err
		}

		if res.StatusCode() == 404 {
			return nil, manga.ErrMangaNotFound
		}

		var batch []paginationEntry
		if err := json.Unmarshal(res.Bytes(), &batch); err != nil {
			return nil, fmt.Errorf("%w: %w", manga.ErrMalformedChapterData, err)
		}

		entries = append(entries, batch...)
		if len(batch) < paginationPageSize {
			break
		}
	}
	return entries, nil
}

const episodeJSONID = `id="episode-json"`
const episodeJSONIDAlt = `id='episode-json'`

func (e *Extractor) fetchEpisodeJSON(ctx context.Context, URL string) (*episodeJSON, error) {
	req := httpClient.R().SetContext(ctx).SetHeader("User-Agent", iPhoneUA)
	e.ApplyCookie(req)

	res, err := req.Get(URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, manga.ErrChapterNotFound
	}

	body := res.String()

	idIdx := strings.Index(body, episodeJSONID)
	quote := byte('"')
	if idIdx == -1 {
		idIdx = strings.Index(body, episodeJSONIDAlt)
		quote = '\''
	}
	if idIdx == -1 {
		return nil, fmt.Errorf("%w: episode-json not found", manga.ErrMalformedChapterData)
	}

	marker := fmt.Sprintf("data-value=%c", quote)
	markerIdx := strings.Index(body[idIdx:], marker)
	if markerIdx == -1 {
		return nil, fmt.Errorf("%w: episode-json data-value not found", manga.ErrMalformedChapterData)
	}
	start := idIdx + markerIdx + len(marker)

	end := strings.IndexByte(body[start:], quote)
	if end == -1 {
		return nil, fmt.Errorf("%w: episode-json unterminated", manga.ErrMalformedChapterData)
	}

	raw := html.UnescapeString(body[start : start+end])

	var ej episodeJSON
	if err := json.Unmarshal([]byte(raw), &ej); err != nil {
		return nil, fmt.Errorf("%w: %w", manga.ErrMalformedChapterData, err)
	}

	return &ej, nil
}

func init() {
	registry.Register("shonenjumpplus.com", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
