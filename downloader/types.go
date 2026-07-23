package downloader

import (
	"context"
	"github.com/knst0/mdl/sdk/manga"
	"sync"
)

type (
	downloadPageFunc func(ctx context.Context, dir string, page *manga.Page) error

	Downloader struct {
		ctx          context.Context
		ch           chan *queueInfo
		wg           sync.WaitGroup
		downloadPage downloadPageFunc
		reporter     ProgressReporter
	}

	queueInfo struct {
		URL       string
		ChapterID string
		Pages     []*manga.Page
	}
)
