package takecomic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // NOTE: register JPEG decoder
	"image/png"
	"regexp"
	"sort"
	"strconv"

	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/sdk/manga"
	"github.com/knst0/mdl/sdk/manga/pluginutil"

	"resty.dev/v3"
)

const (
	scramblePieces = 16
	baseURL        = "https://takecomic.jp"
)

var (
	httpClient = resty.New()
	epURLRe    = regexp.MustCompile(`takecomic\.jp/episodes/([a-f0-9]+)`)
)

type Extractor struct {
	pluginutil.Base
}

func (e *Extractor) FindChapters(ctx context.Context, url string) ([]*manga.Chapter, error) {
	chapter, err := e.FindChapter(ctx, url)
	if err != nil {
		return nil, err
	}
	return []*manga.Chapter{chapter}, nil
}

func (e *Extractor) FindChapter(ctx context.Context, url string) (*manga.Chapter, error) {
	episodeID := extractEpisodeID(url)
	if episodeID == "" {
		return nil, manga.ErrInvalidURLFormat
	}

	ep, err := fetchEpisode(ctx, episodeID)
	if err != nil {
		return nil, err
	}

	viewerID := getViewerID(ep)
	if viewerID == "" {
		return nil, fmt.Errorf("takecomic: no viewerId in episode content: %w", manga.ErrMalformedChapterData)
	}

	chapter := &manga.Chapter{
		ID:      viewerID,
		Number:  strconv.Itoa(ep.SeriesEpisodeNumber),
		Title:   ep.Summary.Title,
		Index:   uint(ep.SeriesEpisodeNumber - 1),
		URL:     url,
		MangaID: ep.Series.ID,
	}
	return chapter, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	viewerID := chapter.ID

	ci, err := fetchContentsInfo(ctx, viewerID)
	if err != nil {
		return nil, err
	}

	pages := ci.Result
	sort.Slice(pages, func(i, j int) bool { return pages[i].Sort < pages[j].Sort })

	return pluginutil.BuildPages(len(pages), ".png", func(index int, filename string) (*manga.Page, error) {
		page := pages[index]
		scramble := parseScramble(page.Scramble)

		return &manga.Page{
			URL:      page.ImageURL,
			Filename: filename,
			Index:    uint(index),
			Headers:  map[string]string{"Referer": baseURL + "/"},
			Decode: func(b []byte) ([]byte, error) {
				if len(scramble) != scramblePieces {
					// No valid scramble — return image as-is
					return b, nil
				}
				return descrambleJPEG(b, scramble)
			},
		}, nil
	})
}

// extractEpisodeID pulls the hex episode hash from a URL like
// https://takecomic.jp/episodes/77e7e45dcb279.
func extractEpisodeID(url string) string {
	m := epURLRe.FindStringSubmatch(url)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

func fetchEpisode(ctx context.Context, episodeID string) (*Episode, error) {
	resp, err := httpClient.R().
		SetContext(ctx).
		Get(baseURL + "/api/episodes/" + episodeID)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == 404 {
		return nil, manga.ErrChapterNotFound
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("takecomic: episode API returned %d", resp.StatusCode())
	}

	var er EpisodeResponse
	if err := json.Unmarshal(resp.Bytes(), &er); err != nil {
		return nil, err
	}
	return &er.Episode, nil
}

func getViewerID(ep *Episode) string {
	for _, c := range ep.Content {
		if c.Type == "viewer" && c.ViewerID != "" {
			return c.ViewerID
		}
	}
	return ""
}

func fetchContentsInfo(ctx context.Context, viewerID string) (*ContentsInfoResponse, error) {
	resp, err := httpClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"user-id":          "",
			"comici-viewer-id": viewerID,
			"page-from":        "0",
			"page-to":          "0",
		}).
		Get(baseURL + "/api/book/contentsInfo")
	if err != nil {
		return nil, err

	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("takecomic: contentsInfo API returned %d", resp.StatusCode())
	}

	var first ContentsInfoResponse
	if err := json.Unmarshal(resp.Bytes(), &first); err != nil {
		return nil, err
	}

	total := first.TotalPages
	if total <= 1 {
		return &first, nil
	}

	resp, err = httpClient.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"user-id":          "",
			"comici-viewer-id": viewerID,
			"page-from":        "0",
			"page-to":          strconv.Itoa(total - 1),
		}).
		Get(baseURL + "/api/book/contentsInfo")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("takecomic: contentsInfo API returned %d", resp.StatusCode())
	}

	var all ContentsInfoResponse
	if err := json.Unmarshal(resp.Bytes(), &all); err != nil {
		return nil, err
	}
	return &all, nil
}

// parseScramble parses a JSON array string like "[6,3,0,14,...]" into an int slice.
func parseScramble(raw string) []int {
	if raw == "" || raw == "[]" {
		return nil
	}
	var arr []int
	if err := json.Unmarshal([]byte(raw), &arr); err != nil {
		return nil
	}
	return arr
}

func descrambleJPEG(raw []byte, scramble []int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("takecomic: decode image: %w", err)
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	const gridCols, gridRows = 4, 4

	type cell struct{ x0, y0, x1, y1 int }
	cells := make([]cell, scramblePieces)
	colW := w / gridCols
	rowH := h / gridRows
	for c := 0; c < gridCols; c++ {
		for r := 0; r < gridRows; r++ {
			idx := c*gridRows + r
			cells[idx] = cell{
				x0: c * colW,
				y0: r * rowH,
				x1: c*colW + colW,
				y1: r*rowH + rowH,
			}
		}
	}

	dstW := colW * gridCols
	dstH := rowH * gridRows
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	for dstPos, srcPos := range scramble {
		src := cells[srcPos]
		d := cells[dstPos]
		for dy := 0; dy < src.y1-src.y0; dy++ {
			for dx := 0; dx < src.x1-src.x0; dx++ {
				dst.Set(d.x0+dx, d.y0+dy, img.At(src.x0+dx, src.y0+dy))
			}
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		return nil, fmt.Errorf("takecomic: encode png: %w", err)
	}
	return buf.Bytes(), nil
}

func init() {
	registry.Register("takecomic.jp", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
