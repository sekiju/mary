package downloader

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/config"
	"github.com/sekiju/mdl/extractor"
	"github.com/sekiju/mdl/sdk/manga"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func (d *Downloader) Queue(URL string) {
	qi := &queueInfo{URL: URL}

	log.Debug().Msgf("New queue item with URL: %s", URL)

	d.wg.Add(1)
	d.ch <- qi
}

func (d *Downloader) Stop() {
	close(d.ch)
	d.wg.Wait()
}

func (d *Downloader) downloadImages(qi *queueInfo) error {
	destination := filepath.Join(config.Params.Output.Directory, qi.ChapterID)

	if _, err := os.Stat(destination); err == nil && config.Params.Output.CleanOnStart {
		if err = os.RemoveAll(destination); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		log.Error().Err(err).Msgf("Failed to create download directory for chapter %s", qi.ChapterID)
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var failedPages int
	semaphore := make(chan struct{}, config.Params.Application.MaxParallelDownloads)

	for _, page := range qi.Pages {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(page *manga.Page) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			if err := d.downloadPage(destination, page); err != nil {
				log.Error().Err(err).Msgf("Failed to download page #%d", page.Index)
				mu.Lock()
				failedPages++
				mu.Unlock()
				return
			}
		}(page)
	}

	wg.Wait()

	if failedPages > 0 {
		return fmt.Errorf("failed to download %d of %d page(s)", failedPages, len(qi.Pages))
	}

	return nil
}

func (d *Downloader) run() {
	semaphore := make(chan struct{}, 1)

	for qi := range d.ch {
		semaphore <- struct{}{}
		go func(qi *queueInfo) {
			defer func() {
				<-semaphore
				d.wg.Done()
			}()

			log.Info().Str("url", qi.URL).Msg("Downloading next chapter in queue")
			start := time.Now()

			parsedURL, err := url.Parse(qi.URL)
			if err != nil {
				log.Error().Str("url", qi.URL).Err(err).Msg("Invalid chapter URL")
				return
			}

			ext, err := extractor.NewExtractor(parsedURL.Hostname())
			if err != nil {
				log.Error().Err(err).Str("url", qi.URL).Send()
				return
			}

			chapter, err := ext.FindChapter(qi.URL)
			if err != nil {
				log.Error().Str("url", qi.URL).Err(err).Msg("Failed to find chapter")
				return
			}

			qi.ChapterID = chapter.ID

			pages, err := ext.FindChapterPages(chapter)
			if err != nil {
				log.Error().Str("chapterId", chapter.ID).Err(err).Msg("Failed to find chapter pages")
				return
			}

			qi.Pages = pages

			if err = d.downloadImages(qi); err != nil {
				log.Error().Str("chapterId", qi.ChapterID).Err(err).Msg("Failed to download chapter")
				return
			}

			log.Info().Str("chapterId", qi.ChapterID).Str("duration", time.Since(start).String()).Msg("Download complete")

		}(qi)
	}
}

func NewDownloader() *Downloader {
	var downloadFunc downloadPageFunc
	if config.Params.Output.FileFormat == config.AutoOutputFormat {
		downloadFunc = func(dir string, page *manga.Page) error {
			r, err := getReader(page)
			if err != nil {
				return err
			}

			return saveFile(dir, page.Filename, r)
		}
	} else {
		downloadFunc = func(dir string, page *manga.Page) error {
			r, err := getReader(page)
			if err != nil {
				return err
			}

			return saveEncodedImage(dir, page.Filename, config.Params.Output.FileFormat, r)
		}
	}

	d := &Downloader{ch: make(chan *queueInfo), downloadPage: downloadFunc}

	go d.run()

	return d
}
