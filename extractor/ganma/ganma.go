package ganma

import (
	"errors"
	"fmt"
	json "github.com/bytedance/sonic"
	util "github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"regexp"
	"resty.dev/v3"
	"strconv"
)

type Extractor struct {
	settings *manga.Settings
}

const (
	sha256hash = "75f44fb799c1d505ae245b52633b59e3db9be1e2fe90b4c766eb5d96a86d5be7"
)

func (e *Extractor) FindChapters(URL string) ([]*manga.Chapter, error) {
	if e.settings.Cookie == nil {
		return nil, manga.ErrCredentialsRequired
	}

	res, err := resty.New().R().Get(URL)
	if err != nil {
		return nil, err
	}

	text := res.String()

	mangaID, err := util.ExtractStringFromHTML(text, `\"magazineId\":\"`, `\"`)
	if err != nil {
		return nil, err
	}

	res, err = resty.New().R().
		SetHeader("Cookie", *e.settings.Cookie).
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
		chapters[i] = &manga.Chapter{
			ID:      item.StoryId,
			Number:  strconv.Itoa(item.Number),
			Title:   item.Title,
			Index:   uint(i),
			URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", magazine.Root.Id, item.StoryId),
			MangaID: magazine.Root.Id,
		}
	}

	return chapters, nil
}

var re = regexp.MustCompile(`https://ganma.jp/web/reader/([a-zA-Z0-9_-]*)/([a-zA-Z0-9_-]*)`)

func (e *Extractor) FindChapter(URL string) (*manga.Chapter, error) {
	matches := re.FindStringSubmatch(URL)
	if len(matches) != 3 {
		return nil, manga.ErrInvalidURLFormat
	}

	if e.settings.Cookie == nil {
		return nil, manga.ErrCredentialsRequired
	}

	var mangaID string
	if util.IsValidUUID(matches[1]) {
		mangaID = matches[1]
	} else {
		res, err := resty.New().R().Get(URL)
		if err != nil {
			return nil, err
		}

		text := res.String()

		mangaID, err = util.ExtractStringFromHTML(text, `\"magazineId\":\"`, `\"`)
		if err != nil {
			return nil, err
		}
	}

	res, err := resty.New().R().
		SetHeader("Cookie", *e.settings.Cookie).
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
			return &manga.Chapter{
				ID:      item.StoryId,
				Number:  strconv.Itoa(item.Number),
				Title:   item.Title,
				Index:   0,
				URL:     fmt.Sprintf("https://ganma.jp/web/reader/%s/%s/0", magazine.Root.Id, item.StoryId),
				MangaID: magazine.Root.Id,
			}, nil
		}
	}

	return nil, manga.ErrChapterNotFound
}

func (e *Extractor) FindChapterPages(chapter *manga.Chapter) ([]*manga.Page, error) {
	res, err := resty.New().R().
		SetHeader("Cookie", *e.settings.Cookie).
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

	pages := make([]*manga.Page, reader.Data.Magazine.StoryContents.PageImages.PageCount)
	padRenamer := renamer.New(reader.Data.Magazine.StoryContents.PageImages.PageCount)

	for index := range reader.Data.Magazine.StoryContents.PageImages.PageCount {
		pages[index] = &manga.Page{
			URL:      fmt.Sprintf("%s%d.jpg?%s&w=4000", reader.Data.Magazine.StoryContents.PageImages.PageImageBaseURL, index+1, reader.Data.Magazine.StoryContents.PageImages.PageImageSign),
			Filename: padRenamer.Name(index, ".jpeg"),
			Index:    uint(index),
		}
	}

	return pages, nil
}

func (e *Extractor) SetSettings(settings manga.Settings) {
	e.settings = &settings
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

func New() (manga.Extractor, error) {
	return &Extractor{settings: &manga.Settings{}}, nil
}
