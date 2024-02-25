package manga_bilibili

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"mary/internal/static"
	"mary/internal/utils"
	"mary/pkg/request"
)

type MangaBiliBili struct {
	domain string
}

func (c *MangaBiliBili) Data() *static.ConnectorData {
	return &static.ConnectorData{
		Domain:               c.domain,
		AuthorizationStatus:  static.AuthorizationStatusOptional,
		ChapterListAvailable: true,
	}
}

func (c *MangaBiliBili) ResolveType(uri url.URL) (static.UrlType, error) {
	pathSegments := strings.Split(uri.Path, "/")

	if len(pathSegments) < 3 {
		return "", errors.New("invalid URL format")
	}

	if pathSegments[1] == "detail" {
		return static.UrlTypeBook, nil
	} else if pathSegments[2] != "" {
		return static.UrlTypeChapter, nil
	}

	return "", errors.New("unable to resolve URL type")
}

func (c *MangaBiliBili) Book(_ url.URL) (*static.Book, error) {
	return nil, static.MassiveDownloaderUnsupportedErr
}

func (c *MangaBiliBili) Chapter(uri url.URL) (*static.Chapter, error) {
	id, err := strconv.Atoi(utils.LastURLSegment(uri.Path))
	if err != nil {
		return nil, err
	}

	res, err := request.Post[GetEpisodeResponse](
		"https://manga.bilibili.com/twirp/comic.v1.Comic/GetEpisode?device=pc&platform=web",
		request.SetHeader("Origin", "https://manga.bilibili.com"),
		request.Body(GetEpisodeRequest{ID: id}),
	)
	if err != nil {
		return nil, err
	}

	return &static.Chapter{
		ID:    id,
		Title: res.Body.Data.ShortTitle,
		Error: nil,
	}, nil
}

func (c *MangaBiliBili) Pages(chapterID any, imageChan chan<- static.Image) error {
	res, err := request.Post[GetImageIndexResponse](
		"https://manga.bilibili.com/twirp/comic.v1.Comic/GetImageIndex?device=pc&platform=web",
		request.SetHeader("Origin", "https://manga.bilibili.com"),
		request.Body(GetImageIndexRequest{EpId: chapterID.(int)}),
	)
	if err != nil {
		return err
	}

	log.Trace().Msgf("data.host: %s", res.Body.Data.Host)

	urls := make([]string, 0)
	for _, image := range res.Body.Data.Images {
		urls = append(urls, image.Path)
	}

	tokenResponse, err := request.Post[ImageTokenResponse](
		"https://manga.bilibili.com/twirp/comic.v1.Comic/ImageToken?device=pc&platform=web",
		request.SetHeader("Origin", "https://manga.bilibili.com"),
		request.Body(ImageTokenRequest{Urls: fmt.Sprintf("[\"%s\"]", strings.Join(urls, "\",\""))}),
	)
	if err != nil {
		return err
	}

	indexNamer := utils.NewIndexNamer(len(tokenResponse.Body.Data))
	for i, datum := range tokenResponse.Body.Data {
		imgUri := fmt.Sprintf("%s?token=%s", datum.Url, datum.Token)
		log.Trace().Msgf("image: %s", imgUri)

		ext := filepath.Ext(path.Base(datum.Url))

		var fn static.ImageFn
		fn = func() ([]byte, error) {
			res, err := request.Get[[]byte](imgUri, request.SetHeader("Origin", "https://manga.bilibili.com"))
			if err != nil {
				return nil, err
			}

			return res.Body, nil
		}

		imageChan <- static.NewImage(indexNamer.Get(i, ext), &fn)
	}

	close(imageChan)
	return nil
}

func New() *MangaBiliBili {
	return &MangaBiliBili{domain: "manga.bilibili.com"}
}
