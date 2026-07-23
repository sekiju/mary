//go:build integration

package giga_viewer

import (
	"context"
	"github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/sdk/manga"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProvider(t *testing.T) {
	ctx := context.Background()
	ext, _ := New()

	t.Run("FindChapters", func(t *testing.T) {
		episodes, err := ext.FindChapters(ctx, "https://shonenjumpplus.com/episode/3269754496608909464")
		assert.NoError(t, err)
		assert.Len(t, episodes, 3)

		_, err = ext.FindChapters(ctx, "https://shonenjumpplus.com/episode/07106371892806261346")
		assert.Equal(t, manga.ErrMangaNotFound, err)
	})

	t.Run("FindChapter", func(t *testing.T) {
		chapter, err := ext.FindChapter(ctx, "https://shonenjumpplus.com/episode/17106371892806261346")
		assert.NoError(t, err)
		assert.Equal(t, "17106371892806261346", chapter.ID)

		_, err = ext.FindChapter(ctx, "https://shonenjumpplus.com/episode/07106371892806261346")
		assert.Equal(t, manga.ErrChapterNotFound, err)

		t.Run("ExtractEpisode", func(t *testing.T) {
			pages, err := ext.FindChapterPages(ctx, chapter)
			if err == manga.ErrPaidChapter {
				t.Skip("test chapter is no longer free to read; needs a currently-free replacement URL")
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, pages)
			util.AssertImage(t, "testdata/giga_viewer__sjp__episode_17106371892806261346.golden", pages[0])
		})
	})
}
