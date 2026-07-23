package downloader

import (
	"testing"

	"github.com/sekiju/mdl/config"
)

func TestChapterConcurrency_FallsBackToOneWhenUnset(t *testing.T) {
	original := config.Params.File.Application.MaxParallelChapters
	defer func() { config.Params.File.Application.MaxParallelChapters = original }()

	config.Params.File.Application.MaxParallelChapters = 0
	if got := chapterConcurrency(); got != 1 {
		t.Errorf("chapterConcurrency() = %d, want 1", got)
	}

	config.Params.File.Application.MaxParallelChapters = -1
	if got := chapterConcurrency(); got != 1 {
		t.Errorf("chapterConcurrency() = %d, want 1", got)
	}
}

func TestChapterConcurrency_UsesConfiguredBound(t *testing.T) {
	original := config.Params.File.Application.MaxParallelChapters
	defer func() { config.Params.File.Application.MaxParallelChapters = original }()

	config.Params.File.Application.MaxParallelChapters = 3
	if got := chapterConcurrency(); got != 3 {
		t.Errorf("chapterConcurrency() = %d, want 3", got)
	}
}
