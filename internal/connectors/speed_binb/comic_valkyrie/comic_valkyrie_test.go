package comic_valkyrie

import (
	"mary/internal/utils"
	"testing"
)

func TestComicValkyrie_Chapter(t *testing.T) {
	utils.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/comic_valkyrie.jpg",
		"https://www.comic-valkyrie.com/samplebook/val_spakuyaku01/",
	)
}
