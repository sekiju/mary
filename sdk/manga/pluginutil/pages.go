package pluginutil

import (
	"github.com/sekiju/mdl/internal/renamer"
	"github.com/sekiju/mdl/sdk/manga"
)

// BuildPages renames count pages with a shared renamer.Renamer and builds
// each *manga.Page via build, stopping and propagating the first error.
func BuildPages(count int, ext string, build func(index int, filename string) (*manga.Page, error)) ([]*manga.Page, error) {
	pages := make([]*manga.Page, count)
	padRenamer := renamer.New(count)

	for i := 0; i < count; i++ {
		page, err := build(i, padRenamer.Name(i, ext))
		if err != nil {
			return nil, err
		}

		pages[i] = page
	}

	return pages, nil
}
