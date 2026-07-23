//go:build integration

package storia_takeshobo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	ext, _ := New()

	t.Run("FindChapter", func(t *testing.T) {
		chapter, err := ext.FindChapter("https://storia.takeshobo.co.jp/_files/mahoako/01/")
		assert.NoError(t, err)
		assert.Equal(t, "01", chapter.ID)
		assert.Equal(t, "mahoako", chapter.MangaID)
	})
}
