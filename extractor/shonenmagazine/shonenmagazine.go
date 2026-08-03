package shonenmagazine

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"image"
	"image/draw"
	"image/jpeg" // NOTE: also registers the JPEG decoder
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

const (
	apiBase     = "https://api.pocket.shonenmagazine.com"
	seAPIBase   = "https://se-api.pocket.shonenmagazine.com"
	platform    = "3"
	isCrawler   = "false"
	listPageLen = 50
)

var httpClient = resty.New()

type Extractor struct {
	pluginutil.Base
}

var (
	titleRe   = regexp.MustCompile(`https://pocket\.shonenmagazine\.com/title/(\d+)`)
	episodeRe = regexp.MustCompile(`https://pocket\.shonenmagazine\.com/title/(\d+)/episode/(\d+)`)
)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := episodeRe.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidChapterURL
	}
	episodeID := matches[2]

	entries, err := e.fetchEpisodeList(ctx, []string{episodeID})
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, manga.ErrChapterNotFound
	}
	ep := entries[0]

	return &manga.Chapter{
		ID:      strconv.FormatInt(ep.EpisodeID, 10),
		Number:  strconv.Itoa(ep.Index),
		Title:   ep.EpisodeName,
		URL:     URL,
		MangaID: strconv.FormatInt(ep.TitleID, 10),
	}, nil
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	vj, err := e.fetchViewer(ctx, chapter.ID)
	if err != nil {
		return nil, err
	}

	titleID := vj.TitleID
	episodeID, err := strconv.ParseInt(chapter.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid episode id %q", manga.ErrMalformedChapterData, chapter.ID)
	}

	if len(vj.PageList) == 0 {
		return nil, fmt.Errorf("%w: no pages found for chapter %s", manga.ErrMalformedChapterData, chapter.ID)
	}

	return pluginutil.BuildPages(len(vj.PageList), ".jpg", func(index int, filename string) (*manga.Page, error) {
		return &manga.Page{
			Index:    uint(index),
			URL:      vj.PageList[index],
			Filename: filename,
			Decode: func(raw []byte) ([]byte, error) {
				return descramblePuzzleImage(raw, titleID, episodeID, vj.ScrambleSeed)
			},
		}, nil
	})
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	titleID := ""

	if matches := titleRe.FindStringSubmatch(URL); len(matches) == 2 {
		titleID = matches[1]
	} else if matches := episodeRe.FindStringSubmatch(URL); len(matches) == 3 {
		entries, err := e.fetchEpisodeList(ctx, []string{matches[2]})
		if err != nil {
			return nil, err
		}
		if len(entries) == 0 {
			return nil, manga.ErrChapterNotFound
		}
		titleID = strconv.FormatInt(entries[0].TitleID, 10)
	} else {
		return nil, manga.ErrInvalidURLFormat
	}

	td, err := e.fetchTitleDetail(ctx, titleID)
	if err != nil {
		return nil, err
	}
	ids := td.WebTitle.EpisodeIDList
	if len(ids) == 0 {
		return nil, manga.ErrChapterNotFound
	}

	entries := make([]episodeEntry, 0, len(ids))
	for start := 0; start < len(ids); start += listPageLen {
		end := min(start+listPageLen, len(ids))
		batch := make([]string, 0, end-start)
		for _, id := range ids[start:end] {
			batch = append(batch, strconv.FormatInt(id, 10))
		}
		es, err := e.fetchEpisodeList(ctx, batch)
		if err != nil {
			return nil, err
		}
		entries = append(entries, es...)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Index < entries[j].Index
	})

	chapters := make([]*manga.Chapter, 0, len(entries))
	for i, ep := range entries {
		chapters = append(chapters, &manga.Chapter{
			ID:      strconv.FormatInt(ep.EpisodeID, 10),
			Number:  strconv.Itoa(ep.Index),
			Title:   ep.EpisodeName,
			Index:   uint(i),
			URL:     episodeURL(ep.TitleID, ep.EpisodeID),
			MangaID: strconv.FormatInt(ep.TitleID, 10),
		})
	}

	return chapters, nil
}

func episodeURL(titleID, episodeID int64) string {
	return fmt.Sprintf("https://pocket.shonenmagazine.com/title/%05d/episode/%d", titleID, episodeID)
}

func (e *Extractor) fetchViewer(ctx context.Context, episodeID string) (*viewerJSON, error) {
	params := map[string]string{"episode_id": episodeID}

	res, err := e.newRequest(ctx).SetQueryParams(params).
		SetHeader("x-manga-hash", serviceHash(params)).
		Get(seAPIBase + "/web/episode/viewer")
	if err != nil {
		return nil, err
	}

	var vj viewerJSON
	if err = json.Unmarshal(res.Bytes(), &vj); err != nil {
		return nil, fmt.Errorf("%w: %w", manga.ErrMalformedChapterData, err)
	}
	if vj.Status != "success" {
		return nil, fmt.Errorf("%w: %s", manga.ErrChapterNotFound, vj.ErrorMessage)
	}
	return &vj, nil
}

func (e *Extractor) fetchEpisodeList(ctx context.Context, ids []string) ([]episodeEntry, error) {
	params := map[string]string{"episode_id_list": strings.Join(ids, ",")}

	res, err := e.newRequest(ctx).SetFormData(params).
		SetHeader("x-manga-hash", serviceHash(params)).
		Post(apiBase + "/episode/list")
	if err != nil {
		return nil, err
	}

	var el episodeListJSON
	if err := json.Unmarshal(res.Bytes(), &el); err != nil {
		return nil, fmt.Errorf("%w: %w", manga.ErrMalformedChapterData, err)
	}
	if el.Status != "success" {
		return nil, fmt.Errorf("%w: %s", manga.ErrChapterNotFound, el.ErrorMessage)
	}
	return el.EpisodeList, nil
}

func (e *Extractor) fetchTitleDetail(ctx context.Context, titleID string) (*titleDetailJSON, error) {
	params := map[string]string{"title_id": titleID}

	res, err := e.newRequest(ctx).SetQueryParams(params).
		SetHeader("x-manga-hash", serviceHash(params)).
		Get(apiBase + "/web/title/detail")
	if err != nil {
		return nil, err
	}

	var td titleDetailJSON
	if err = json.Unmarshal(res.Bytes(), &td); err != nil {
		return nil, fmt.Errorf("%w: %w", manga.ErrMalformedChapterData, err)
	}
	if td.Status != "success" {
		return nil, fmt.Errorf("%w: %s", manga.ErrMangaNotFound, td.ErrorMessage)
	}
	return &td, nil
}

func (e *Extractor) newRequest(ctx context.Context) *resty.Request {
	req := httpClient.R().SetContext(ctx).
		SetHeader("x-manga-platform", platform).
		SetHeader("x-manga-is-crawler", isCrawler).
		SetHeader("Referer", "https://pocket.shonenmagazine.com/")
	e.ApplyCookie(req)
	return req
}

// serviceHash signs the request params: every sorted key/value pair is
// hashed (sha256(key)_sha512(value)), the pairs joined with commas are
// hashed again, and the result is finalized with sha512 over the empty
// key/value pair. The server rejects unsigned requests (response 1001).
func serviceHash(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, digest(sha256.New(), k)+"_"+digest(sha512.New(), params[k]))
	}

	joined := digest(sha256.New(), strings.Join(parts, ","))
	emptyPair := digest(sha256.New(), "") + "_" + digest(sha512.New(), "")
	return digest(sha512.New(), joined+emptyPair)
}

func digest(h hash.Hash, s string) string {
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	scrambleDivideNum = 4
	scrambleMultiple  = 8
)

func descramblePuzzleImage(raw []byte, titleID, episodeID int64, baseSeed string) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("shonenmagazine: decode image: %w", err)
	}

	order, err := scrambleOrder(titleID, episodeID, baseSeed)
	if err != nil {
		return nil, fmt.Errorf("shonenmagazine: %w", err)
	}

	bounds := img.Bounds()
	cellW := (bounds.Dx() / (scrambleDivideNum * scrambleMultiple)) * scrambleMultiple
	cellH := (bounds.Dy() / (scrambleDivideNum * scrambleMultiple)) * scrambleMultiple

	dst := image.NewRGBA(bounds)
	for dstIdx, srcIdx := range order {
		dstX, dstY := (dstIdx%scrambleDivideNum)*cellW, (dstIdx/scrambleDivideNum)*cellH
		srcX, srcY := (srcIdx%scrambleDivideNum)*cellW, (srcIdx/scrambleDivideNum)*cellH
		draw.Draw(dst, image.Rect(dstX, dstY, dstX+cellW, dstY+cellH),
			img, image.Pt(srcX, srcY), draw.Src)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: 90}); err != nil {
		return nil, fmt.Errorf("shonenmagazine: encode jpeg: %w", err)
	}
	return buf.Bytes(), nil
}

// scrambleOrder maps each 4x4 destination tile to the source tile that
// belongs there, derived from the episode's scramble seed. The seed
// characters are indices into a title-parity alphabet, XORed with the
// title/episode ids, then fed through a xorshift32 PRNG whose 16 outputs
// are sorted to produce the permutation.
func scrambleOrder(titleID, episodeID int64, baseSeed string) ([scrambleDivideNum * scrambleDivideNum]int, error) {
	var order [scrambleDivideNum * scrambleDivideNum]int

	alphabet := "q6jtf2xnog"
	if titleID%2 == 0 {
		alphabet = "svdk0m7acl"
	}

	var sb strings.Builder
	for _, c := range baseSeed {
		idx := strings.IndexRune(alphabet, c)
		if idx == -1 {
			return order, fmt.Errorf("invalid scramble seed %q", baseSeed)
		}
		sb.WriteByte(byte('0' + idx))
	}

	seed, err := strconv.ParseUint(sb.String(), 10, 64)
	if err != nil {
		return order, fmt.Errorf("invalid scramble seed %q", baseSeed)
	}

	x := uint32(seed) ^ uint32(titleID+episodeID)

	type prngPair struct {
		val, idx uint32
	}
	pairs := make([]prngPair, 0, scrambleDivideNum*scrambleDivideNum)
	for i := range scrambleDivideNum * scrambleDivideNum {
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
		pairs = append(pairs, prngPair{x, uint32(i)})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })

	for i, p := range pairs {
		order[i] = int(p.idx)
	}
	return order, nil
}

func init() {
	registry.Register("pocket.shonenmagazine.com", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
