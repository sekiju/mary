package takecomic

import (
	"context"
	"testing"

	"github.com/knst0/mdl/extractor/testutil"
)

func TestChapterSnapshot_E2E(t *testing.T) {
	const chapterURL = "https://takecomic.jp/episodes/77e7e45dcb279"

	e, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
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
