package ganma

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	json "github.com/bytedance/sonic"
	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/sdk/manga"

	"resty.dev/v3"
)

type Extractor struct{}

func (e *Extractor) SetSettings(_ manga.Settings) {}

func toChapterSimple(storyID, title, magazineID string) manga.Chapter {
	return manga.Chapter{
		ID:      storyID,
		Title:   title,
		URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", magazineID, storyID),
		MangaID: magazineID,
	}
}

// readerPageProps holds extracted reader page component props
type readerPageProps struct {
	Title                  string `json:"title"`
	StoryTitle             string `json:"storyTitle"`
	MagazineTitle          string `json:"magazineTitle"`
	MagazineID             string `json:"magazineId"`
	MagazineIDOrAlias      string `json:"magazineIdOrAlias"`
	StoryID                string `json:"storyId"`
	StoryPageCount         int    `json:"storyPageCount"`
	SingleModeDisplayUnits []struct {
		URL  string `json:"url"`
		Page int    `json:"page"`
	} `json:"singleModeDisplayUnits"`
}

func extractReaderProps(html string) (*readerPageProps, error) {
	const marker = `\"singleModeDisplayUnits\":[`

	idx := strings.Index(html, marker)
	if idx == -1 {
		return nil, errors.New("reader props marker not found in HTML")
	}

	nullBrace := strings.LastIndex(html[:idx], `null,{`)
	if nullBrace == -1 {
		return nil, errors.New("props object start not found")
	}
	start := nullBrace + 5

	depth := 0
	end := -1
	for i := start; i < len(html); i++ {
		switch html[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i + 1
				goto done
			}
		}
	}
done:
	if end == -1 {
		return nil, errors.New("props object end not found")
	}

	raw := html[start:end]
	unescaped := strings.ReplaceAll(raw, `\"`, `"`)
	unescaped = strings.ReplaceAll(unescaped, `\\`, `\`)
	unescaped = strings.ReplaceAll(unescaped, `\n`, "\n")

	var props readerPageProps
	if err := json.Unmarshal([]byte(unescaped), &props); err != nil {
		return nil, fmt.Errorf("failed to parse reader props: %w", err)
	}

	return &props, nil
}

var urlRe = regexp.MustCompile(`https://ganma.jp/web/reader/([a-zA-Z0-9_-]+)/([a-zA-Z0-9_-]+)`)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := urlRe.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	res, err := resty.New().R().SetContext(ctx).Get(URL)
	if err != nil {
		return nil, err
	}

	props, err := extractReaderProps(res.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", manga.ErrChapterNotFound, err)
	}

	chapter := toChapterSimple(props.StoryID, props.Title, props.MagazineIDOrAlias)
	return &chapter, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := resty.New().R().SetContext(ctx).Get(chapter.URL)
	if err != nil {
		return nil, err
	}

	props, err := extractReaderProps(res.String())
	if err != nil {
		return nil, fmt.Errorf("failed to extract page data: %w", err)
	}

	pages := make([]*manga.Page, 0, len(props.SingleModeDisplayUnits))
	for _, unit := range props.SingleModeDisplayUnits {
		filename := fmt.Sprintf("%03d.jpeg", unit.Page+1)
		pages = append(pages, &manga.Page{
			URL:      unit.URL,
			Filename: filename,
			Index:    uint(unit.Page),
		})
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("%w: no pages found for chapter %s", manga.ErrMalformedChapterData, chapter.ID)
	}

	return pages, nil
}

var magazineRe = regexp.MustCompile(`https://ganma.jp/web/magazine/([a-zA-Z0-9_-]+)`)

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	matches := magazineRe.FindStringSubmatch(URL)
	if len(matches) != 2 {
		return nil, manga.ErrInvalidURLFormat
	}
	alias := matches[1]

	res, err := resty.New().R().SetContext(ctx).
		SetHeader("RSC", "1").
		SetHeader("Next-Router-State-Tree", "%5B%22%22%2C%7B%22children%22%3A%5B%22magazine%22%2C%7B%22children%22%3A%5B%22"+alias+"%22%2C%7B%22children%22%3A%5B%22__PAGE__%22%2C%7B%7D%5D%7D%5D%7D%5D%7D%5D").
		Get(fmt.Sprintf("https://ganma.jp/web/magazine/%s/episodes?_rsc=mdl", alias))
	if err != nil {
		return nil, err
	}

	chapters := parseEpisodeRSC(res.String(), alias)
	if len(chapters) == 0 {
		return nil, manga.ErrChapterNotFound
	}

	return chapters, nil
}

// parseEpisodeRSC extracts story entries from a Next.js RSC flight stream.
// The RSC payload contains alternating storyId and title fields:
//
//	"title":"<magazine title>"
//	"storyId":"<uuid>","title":"第1話"
//	"storyId":"<uuid>","title":"第2話"
//
// The first title is the magazine title (no preceding storyId); skip it.
func parseEpisodeRSC(rscData, alias string) []*manga.Chapter {
	var chapters []*manga.Chapter

	storyRe := regexp.MustCompile(`"storyId":"([a-f0-9-]+)","title":"([^"]+)"`)
	matches := storyRe.FindAllStringSubmatch(rscData, -1)

	for i, m := range matches {
		chapters = append(chapters, &manga.Chapter{
			ID:      m[1],
			Number:  strconv.Itoa(i + 1),
			Title:   m[2],
			URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", alias, m[1]),
			MangaID: alias,
			Index:   uint(i),
		})
	}

	return chapters
}

func init() {
	registry.Register("ganma.jp", func(cookie *string) (manga.Extractor, error) {
		return &Extractor{}, nil
	})
}

func New() (manga.Extractor, error) {
	return &Extractor{}, nil
}
