package comic_valkyrie

import (
	"559/internal/utils/testing_utils"
	"testing"
)

func TestComicValkyrie_Chapter(t *testing.T) {
	testing_utils.Connector(
		t,
		New(),
		"https://host.kireyev.org/559/comic_valkyrie.jpg",
		"https://www.comic-valkyrie.com/samplebook/val_spakuyaku01/",
	)
}
