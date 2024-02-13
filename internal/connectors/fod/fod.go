package fod

import (
	"559/internal/config"
	"559/internal/static"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"github.com/rs/zerolog/log"
	"image"
	"net/url"
	"strconv"
)

type Fod struct {
	domain string
}

func (c *Fod) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusRequired,
		ChapterListAvailable: false,
	}
}

func (c *Fod) ResolveType(_ url.URL) (static.UrlType, error) {
	return static.UrlTypeChapter, nil
}

func (c *Fod) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *Fod) Chapter(uri url.URL) (*static.Chapter, error) {
	bookCredentials, err := extractValuesFromURL(uri)
	if err != nil {
		return nil, err
	}

	res, err := request.Post[DetailResponse]("https://manga.fod.fujitv.co.jp/api/books/detail", request.Body(&bookCredentials), c.requestHeaders())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch episode: %s", err)
	}

	chapter := static.Chapter{
		ID:    bookCredentials,
		Title: res.Body.BookDetail.BookName,
	}

	_, sessionExists := config.Data.Sites[c.domain]
	if sessionExists && res.Body.BookDetail.IsFree && !res.Body.BookDetail.IsPurchased {
		purchaseResponse, err := request.Post[ServerStatusResponse]("https://manga.fod.fujitv.co.jp/api/purchase/buy", request.Body(map[string]any{
			"buy_type": 1,
			"episodes": []map[string]any{
				{
					"episode_id":       &bookCredentials.EpisodeID,
					"discounted_price": 0,
					"cashback_point":   0,
				},
			},
		}), c.requestHeaders())
		if err != nil {
			return nil, fmt.Errorf("failed to buy episode: %s", err)
		}

		if purchaseResponse.Body.ServerStatus.ResultCode != 0 {
			return nil, fmt.Errorf("failed to buy episode, result code: %d", purchaseResponse.Body.ServerStatus.ResultCode)
		}
	}

	return &chapter, nil
}

func (c *Fod) Pages(chapterID any, imageChan chan<- static.Image) error {
	bookCredentials := chapterID.(*BookCredentialsRequest)

	res, err := request.Post[LicenceKeyResponse]("https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser", request.Body(bookCredentials), c.requestHeaders())
	if err != nil {
		return fmt.Errorf("failed to fetch episode: %s", err)
	}

	indexNamer := utils.NewIndexNamer(res.Body.GuardianInfoForBrowser.BookData.PageCount)
	for i := 1; i <= res.Body.GuardianInfoForBrowser.BookData.PageCount; i++ {
		imageUrl, err := cleanURL(res.Body.GuardianInfoForBrowser.GUARDIANSERVER + normalizeUrl(res.Body.GuardianInfoForBrowser.BookData.S3Key) + strconv.Itoa(i) + ".jpg?" + res.Body.GuardianInfoForBrowser.ADDITIONALQUERYSTRING)
		if err != nil {
			return err
		}

		log.Trace().Msgf("image url: %s", imageUrl)

		processPage(imageUrl, res.Body.GuardianInfoAll.KeysForBrowser[i-1], indexNamer.Get(i, ".jpg"), imageChan)
	}

	close(imageChan)
	return nil
}

func processPage(uri, key, fileName string, imageChan chan<- static.Image) {
	var fn static.ImageFn
	fn = func() ([]byte, error) {
		res, err := request.Get[image.Image](uri)
		if err != nil {
			return nil, err
		}

		return descrambleImage(res.Body, key), nil
	}

	imageChan <- static.NewImage(fileName, &fn)
}

func (c *Fod) requestHeaders() request.OptsFn {
	connectorConfig, exists := config.Data.Sites[c.domain]
	return func(c *request.Config) {
		c.Headers["zk-app-version"] = "1.1.26"
		c.Headers["zk-os-type"] = "1"
		c.Headers["zk-safe-search"] = "0"
		if exists {
			c.Headers["zk-session-key"] = connectorConfig.Session
		}
	}
}

func New() *Fod {
	return &Fod{
		domain: "manga.fod.fujitv.co.jp",
	}
}
