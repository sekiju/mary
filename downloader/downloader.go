package downloader

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/knst0/mdl/config"
	"github.com/knst0/mdl/extractor"
	"github.com/knst0/mdl/sdk/manga"
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

func (d *Downloader) downloadImages(ctx context.Context, qi *queueInfo) error {
	destination := filepath.Join(config.Params.File.Output.Directory, qi.ChapterID)

	if _, err := os.Stat(destination); err == nil && config.Params.File.Output.CleanOnStart {
		if err = os.RemoveAll(destination); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(destination, os.ModePerm); err != nil {
		d.reporter.ChapterError(qi.URL, qi.ChapterID, "Failed to create download directory", err)
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var failedPages int
	semaphore := make(chan struct{}, config.Params.File.Application.MaxParallelDownloads)

	for _, page := range qi.Pages {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(page *manga.Page) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			if ctx.Err() != nil {
				return
			}

			err := d.downloadPage(ctx, destination, page)
			d.reporter.PageDownloaded(qi.ChapterID, page.Index, len(qi.Pages), err)
			if err != nil {
				mu.Lock()
				failedPages++
				mu.Unlock()
				return
			}
		}(page)
	}

	wg.Wait()

	if ctx.Err() != nil {
		if err := os.RemoveAll(destination); err != nil {
			d.reporter.ChapterError(qi.URL, qi.ChapterID, "Failed to clean up partial download", err)
		}
		return ctx.Err()
	}

	if failedPages > 0 {
		return fmt.Errorf("failed to download %d of %d page(s)", failedPages, len(qi.Pages))
	}

	return nil
}

func chapterConcurrency() int {
	if n := config.Params.File.Application.MaxParallelChapters; n > 0 {
		return n
	}
	return 1
}

func (d *Downloader) run() {
	semaphore := make(chan struct{}, chapterConcurrency())

	for qi := range d.ch {
		semaphore <- struct{}{}
		go func(qi *queueInfo) {
			defer func() {
				<-semaphore
				d.wg.Done()
			}()

			d.reporter.ChapterStarted(qi.URL)
			start := time.Now()

			parsedURL, err := url.Parse(qi.URL)
			if err != nil {
				d.reporter.ChapterError(qi.URL, "", "Invalid chapter URL", err)
				return
			}

			ext, err := extractor.NewExtractor(parsedURL.Hostname())
			if err != nil {
				d.reporter.ChapterError(qi.URL, "", "Unsupported website", err)
				return
			}

			chapter, err := ext.FindChapter(d.ctx, qi.URL)
			if err != nil {
				d.reporter.ChapterError(qi.URL, "", "Failed to find chapter", err)
				return
			}

			qi.ChapterID = chapter.ID

			pages, err := ext.FindChapterPages(d.ctx, chapter)
			if err != nil {
				d.reporter.ChapterError(qi.URL, chapter.ID, "Failed to find chapter pages", err)
				return
			}

			qi.Pages = pages

			if err = d.downloadImages(d.ctx, qi); err != nil {
				d.reporter.ChapterError(qi.URL, qi.ChapterID, "Failed to download chapter", err)
				return
			}

			d.reporter.ChapterDone(qi.ChapterID, time.Since(start))

		}(qi)
	}
}

func NewDownloader(ctx context.Context, reporter ...ProgressReporter) *Downloader {
	var downloadFunc downloadPageFunc
	if config.Params.File.Output.FileFormat == config.AutoOutputFormat {
		downloadFunc = func(ctx context.Context, dir string, page *manga.Page) error {
			r, err := getReader(ctx, page)
			if err != nil {
				return err
			}

			return saveFile(dir, page.Filename, r)
		}
	} else {
		downloadFunc = func(ctx context.Context, dir string, page *manga.Page) error {
			r, err := getReader(ctx, page)
			if err != nil {
				return err
			}

			return saveEncodedImage(dir, page.Filename, config.Params.File.Output.FileFormat, r)
		}
	}

	var r ProgressReporter = defaultProgressReporter{}
	if len(reporter) > 0 && reporter[0] != nil {
		r = reporter[0]
	}

	d := &Downloader{ctx: ctx, ch: make(chan *queueInfo), downloadPage: downloadFunc, reporter: r}

	go d.run()

	return d
}
