package fod

import (
	"559/internal/config"
	"559/internal/connectors"
	"559/internal/utils"
	"559/pkg/request"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
	"strings"
)

type Fod struct {
	*connectors.Base
}

func New() *Fod {
	return &Fod{
		Base: connectors.NewBase("manga.fod.fujitv.co.jp"),
	}
}

func (f *Fod) Context() *connectors.Base {
	return f.Base
}

func (f *Fod) Pages(uri url.URL, imageChan chan<- connectors.ReaderImage) error {
	bookId, episodeId, err := extractValuesFromURL(uri)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"zk-app-version": "1.1.26",
		"zk-os-type":     "1",
		"zk-safe-search": "0",
	}

	connectorConfig, exists := config.Data.Sites[f.Domain]

	if exists {
		headers["zk-session-key"] = connectorConfig.Session
	}

	log.Trace().Msgf("zk-session-key: %v", connectorConfig.Session)

	resp, err := request.Post[LicenceKeyResponse]("https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser", &request.Config{
		Headers: headers,
		Body: map[string]string{
			"book_id":    bookId,
			"episode_id": episodeId,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to fetch episode: %s", err)
	}

	log.Trace().Msgf("licenceKeyForBrowser: %v", resp)

	if exists && connectorConfig.PurchaseFreeBooks {
		isFullVersion := strings.Contains(resp.GuardianInfoForBrowser.BookData.S3Key, "_001")
		if !isFullVersion {
			_, err := request.Post[interface{}]("https://manga.fod.fujitv.co.jp/api/purchase/buy", &request.Config{
				Headers: headers,
				Body: map[string]any{
					"buy_type": 1,
					"episodes": []map[string]any{
						{
							"episode_id":       episodeId,
							"discounted_price": 0,
							"cashback_point":   0,
						},
					},
				},
			})

			// todo: normalnaya validation dly oshibok

			if err != nil {
				log.Error().Msg("failed to purchase book")
			} else {
				log.Info().Msg("purchase success")
			}

			resp, err = request.Post[LicenceKeyResponse]("https://manga.fod.fujitv.co.jp/api/books/licenceKeyForBrowser", &request.Config{
				Headers: headers,
				Body: map[string]string{
					"book_id":    bookId,
					"episode_id": episodeId,
				},
			})
			if err != nil {
				return fmt.Errorf("failed to fetch episode: %s", err)
			}
		}
	}

	fnf := utils.NewIndexNameFormatter(resp.GuardianInfoForBrowser.BookData.PageCount)
	for i := 1; i <= resp.GuardianInfoAll.DataForBrowser.PageCount; i++ {
		imageUrl, err := cleanURL(resp.GuardianInfoForBrowser.GUARDIANSERVER + normalizeUrl(resp.GuardianInfoForBrowser.BookData.S3Key) + strconv.Itoa(i) + ".jpg?" + resp.GuardianInfoForBrowser.ADDITIONALQUERYSTRING)
		if err != nil {
			return err
		}

		processPage(imageUrl, resp.GuardianInfoAll.KeysForBrowser[i-1], fnf.GetName(i, ".jpg"), imageChan)
		if config.Data.Settings.Debug.Enable {
			processOriginalPage(imageUrl, fnf.GetName(i, "_original.jpg"), imageChan)
		}
	}

	if config.Data.Settings.Debug.Enable {
		imageUrl, err := cleanURL(resp.GuardianInfoForBrowser.GUARDIANSERVER + normalizeUrl(resp.GuardianInfoForBrowser.BookData.S3Key) + strconv.Itoa(1) + ".jpg?" + resp.GuardianInfoForBrowser.ADDITIONALQUERYSTRING)
		if err != nil {
			return err
		}

		err = processKeys(imageUrl, resp.GuardianInfoAll.KeysForBrowser)
		if err != nil {
			return err
		}
	}

	close(imageChan)
	return nil
}
