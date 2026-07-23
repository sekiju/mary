package downloader

import (
	"time"

	"github.com/rs/zerolog/log"
)

// ProgressReporter receives progress and error events for chapter/page
// downloads, decoupling downloader from any specific logging/UI backend.
type ProgressReporter interface {
	ChapterStarted(url string)
	ChapterTitle(chapterID, title string)
	ChapterError(url, chapterID, msg string, err error)
	PageDownloaded(chapterID string, index uint, total int, err error)
	ChapterDone(chapterID string, duration time.Duration)
}

// defaultProgressReporter reproduces the zerolog output the CLI printed
// before ProgressReporter existed.
type defaultProgressReporter struct{}

func (defaultProgressReporter) ChapterStarted(url string) {
	log.Info().Str("url", url).Msg("Downloading next chapter in queue")
}

func (defaultProgressReporter) ChapterTitle(chapterID, title string) {
	log.Info().Str("chapterId", chapterID).Str("title", title).Msg("Chapter title")
}

func (defaultProgressReporter) ChapterError(url, chapterID, msg string, err error) {
	evt := log.Error().Err(err)
	if chapterID != "" {
		evt = evt.Str("chapterId", chapterID)
	} else {
		evt = evt.Str("url", url)
	}
	evt.Msg(msg)
}

func (defaultProgressReporter) PageDownloaded(chapterID string, index uint, total int, err error) {
	if err != nil {
		log.Error().Err(err).Msgf("Failed to download page #%d", index)
	}
}

func (defaultProgressReporter) ChapterDone(chapterID string, duration time.Duration) {
	log.Info().Str("chapterId", chapterID).Str("duration", duration.String()).Msg("Download complete")
}
