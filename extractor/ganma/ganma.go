package ganma

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"

	json "github.com/bytedance/sonic"
	"github.com/knst0/mdl/extractor/registry"
	util "github.com/knst0/mdl/extractor/util"
	"github.com/knst0/mdl/sdk/manga"
	"github.com/knst0/mdl/sdk/manga/pluginutil"

	"resty.dev/v3"
)

type Extractor struct {
	pluginutil.Base
}

const (
	sha256hash = "75f44fb799c1d505ae245b52633b59e3db9be1e2fe90b4c766eb5d96a86d5be7"
)

func toChapter(storyID, title, magazineID string, number int, index uint) manga.Chapter {
	return manga.Chapter{
		ID:      storyID,
		Number:  strconv.Itoa(number),
		Title:   title,
		Index:   index,
		URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", magazineID, storyID),
		MangaID: magazineID,
	}
}

func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
	if e.Settings.Cookie == nil {
		return nil, manga.ErrCredentialsRequired
	}

	res, err := resty.New().R().SetContext(ctx).Get(URL)
	if err != nil {
		return nil, err
	}

	text := res.String()

	mangaID, err := util.ExtractStringFromHTML(text, `\"magazineId\":\"`, `\"`)
	if err != nil {
		return nil, err
	}

	res, err = resty.New().R().SetContext(ctx).
		SetHeader("Cookie", *e.Settings.Cookie).
		SetHeader("X-From", "https://reader.ganma.jp/api/").
		Get(fmt.Sprintf("https://reader.ganma.jp/api/3.2/magazines/%s", mangaID))
	if err != nil {
		return nil, err
	}

	var magazine magazineResult
	if err = json.Unmarshal(res.Bytes(), &magazine); err != nil {
		return nil, err
	}

	chapters := make([]*manga.Chapter, len(magazine.Root.Items))
	for i, item := range magazine.Root.Items {
		chapter := toChapter(item.StoryId, item.Title, magazine.Root.Id, item.Number, uint(i))
		chapters[i] = &chapter
	}

	return chapters, nil
}

var re = regexp.MustCompile(`https://ganma.jp/web/reader/([a-zA-Z0-9_-]*)/([a-zA-Z0-9_-]*)`)

func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	if e.Settings.Cookie == nil {
		return nil, manga.ErrCredentialsRequired
	}

	var mangaID string
	if util.IsValidUUID(matches[1]) {
		mangaID = matches[1]
	} else {
		res, err := resty.New().R().SetContext(ctx).Get(URL)
		if err != nil {
			return nil, err
		}

		text := res.String()

		mangaID, err = util.ExtractStringFromHTML(text, `\"magazineId\":\"`, `\"`)
		if err != nil {
			return nil, err
		}
	}

	res, err := resty.New().R().SetContext(ctx).
		SetHeader("Cookie", *e.Settings.Cookie).
		SetHeader("X-From", "https://reader.ganma.jp/api/").
		Get(fmt.Sprintf("https://reader.ganma.jp/api/3.2/magazines/%s", mangaID))
	if err != nil {
		return nil, err
	}

	var magazine magazineResult
	if err = json.Unmarshal(res.Bytes(), &magazine); err != nil {
		return nil, err
	}

	for _, item := range magazine.Root.Items {
		if item.StoryId == matches[2] {
			chapter := toChapter(item.StoryId, item.Title, magazine.Root.Id, item.Number, 0)
			return &chapter, nil
		}
	}

	return nil, manga.ErrChapterNotFound
}

func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := resty.New().R().SetContext(ctx).
		SetHeader("Cookie", *e.Settings.Cookie).
		SetHeader("X-From", "https://reader.ganma.jp/api/").
		Get(fmt.Sprintf(
			"https://ganma.jp/api/graphql?operationName=MagazineStoryReaderQuery&variables=%s&extensions=%s",
			url.QueryEscape(fmt.Sprintf(`{"magazineIdOrAlias":%q,"storyId":%q,"publicKey":null}`, chapter.MangaID, chapter.ID)),
			url.QueryEscape(fmt.Sprintf(`{"persistedQuery":{"version":1,"sha256Hash":%q}}`, sha256hash)),
		))
	if err != nil {
		return nil, err
	}

	var reader readerResult
	if err = json.Unmarshal(res.Bytes(), &reader); err != nil {
		return nil, err
	}

	for _, gqlErr := range reader.Errors {
		if gqlErr.Extensions.Code == "STORY_COUNT_LIMITED" {
			return nil, manga.ErrPaidChapter
		}
	}

	return pluginutil.BuildPages(reader.Data.Magazine.StoryContents.PageImages.PageCount, ".jpeg", func(index int, filename string) (*manga.Page, error) {
		return &manga.Page{
			URL:      fmt.Sprintf("%s%d.jpg?%s&w=4000", reader.Data.Magazine.StoryContents.PageImages.PageImageBaseURL, index+1, reader.Data.Magazine.StoryContents.PageImages.PageImageSign),
			Filename: filename,
			Index:    uint(index),
		}, nil
	})
}

func (e *Extractor) GenerateCookie() (string, error) {
	res, err := resty.New().R().SetHeader("X-From", "https://reader.ganma.jp/api/").Post("https://reader.ganma.jp/api/1.0/account")
	if err != nil {
		return "", err
	}

	if res.StatusCode() != 200 {
		return "", errors.New("failed to create account")
	}

	var createAccount createAccountResponse
	if err = json.Unmarshal(res.Bytes(), &createAccount); err != nil {
		return "", err
	}

	res, err = resty.New().R().SetHeader("X-From", "https://reader.ganma.jp/api/").
		SetBody(createAccount.Root).
		Post("https://reader.ganma.jp/api/3.0/session")
	if err != nil {
		return "", err
	}

	if res.StatusCode() != 200 {
		return "", errors.New("failed to login with generated account")
	}

	return res.Header().Get("Set-Cookie"), nil
}

func init() {
	registry.Register("ganma.jp", registry.WithSession(New))
}

func New() (manga.Extractor, error) {
	return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
