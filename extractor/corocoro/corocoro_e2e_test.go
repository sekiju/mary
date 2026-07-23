package corocoro

import (
	"context"
	"testing"

	"github.com/knst0/mdl/extractor/testutil"
	"github.com/knst0/mdl/sdk/manga"
)

func TestChapterSnapshot_E2E(t *testing.T) {
	const chapterURL = "https://www.corocoro.jp/chapter/47950/viewer"

	e, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e.SetSettings(manga.Settings{})
	ctx := context.Background()

	chapter, err := e.FindChapter(ctx, chapterURL)
	if err != nil {
		t.Fatalf("FindChapter: %v", err)
	}

	pages, err := e.FindChapterPages(ctx, chapter)
	if err != nil {
		t.Fatalf("FindChapterPages: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("FindChapterPages returned no pages")
	}

	got := testutil.HashPages(t, pages)
	testutil.CompareSnapshot(t, "testdata/chapter.snapshot.json", got)
}
