//go:build integration

package comic_walker

import (
	"context"
	"github.com/sekiju/mdl/extractor/util"
	"github.com/sekiju/mdl/sdk/manga"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	ctx := context.Background()
	ext, _ := New()

	t.Run("FindChapters", func(t *testing.T) {
		episodes, err := ext.FindChapters(ctx, "https://comic-walker.com/detail/KC_005558_S/episodes/KC_0055580000200011_E?episodeType=first")
		assert.NoError(t, err)
		assert.NotEmpty(t, episodes)
		assert.Equal(t, "018f84b1-1d0b-7557-b1e2-7ec22323c494", episodes[0].ID)

		episodes, err = ext.FindChapters(ctx, "https://comic-walker.com/detail/KC_005558_S?episodeType=first")
		assert.NoError(t, err)
		assert.NotEmpty(t, episodes)
		assert.Equal(t, "018f84b1-1d0b-7557-b1e2-7ec22323c494", episodes[0].ID)

		_, err = ext.FindChapters(ctx, "https://comic-walker.com/detail/KC_005558_D?episodeType=first")
		assert.Equal(t, manga.ErrMangaNotFound, err)
	})

	t.Run("FindChapter", func(t *testing.T) {
		chapter, err := ext.FindChapter(ctx, "https://comic-walker.com/detail/KC_005558_S/episodes/KC_0055580000200011_E?episodeType=first")
		assert.NoError(t, err)
		assert.Equal(t, "018f84b1-1d0b-7557-b1e2-7ec22323c494", chapter.ID)

		chapter, err = ext.FindChapter(ctx, "https://comic-walker.com/detail/KC_005558_S?episodeType=first")
		assert.NoError(t, err)
		assert.Equal(t, "018f84b1-1d0b-7557-b1e2-7ec22323c494", chapter.ID)

		_, err = ext.FindChapter(ctx, "https://comic-walker.com/detail/KC_005558_S/episodes/KC_0015580000200011_E?episodeType=first")
		assert.Equal(t, manga.ErrChapterNotFound, err)

		t.Run("ExtractEpisode", func(t *testing.T) {
			pages, err := ext.FindChapterPages(ctx, chapter)
			assert.NoError(t, err)
			assert.NotEmpty(t, pages)
			util.AssertImage(t, "testdata/comic_walker__KC_005558_S_KC_0055580000200011_E.golden", pages[0])
		})
	})
}
