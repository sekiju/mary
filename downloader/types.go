package downloader

import (
	"context"
	"sync"

	"github.com/knst0/mdl/sdk/manga"
)

type (
	downloadPageFunc func(ctx context.Context, dir string, page *manga.Page) error

	Downloader struct {
		ctx          context.Context
		ch           chan *queueInfo
		wg           sync.WaitGroup
		stopOnce     sync.Once
		downloadPage downloadPageFunc
		reporter     ProgressReporter
	}

	queueInfo struct {
		URL       string
		ChapterID string
		Title     string
		Pages     []*manga.Page
	}
)
