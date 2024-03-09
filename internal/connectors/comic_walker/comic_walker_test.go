package comic_walker

import (
	"mary/internal/utils"
	"testing"
)

func TestComicWalker_Chapter(t *testing.T) {
	utils.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/comic_walker.webp",
		"https://comic-walker.com/detail/KC_004526_S/episodes/KC_0045260000100011_E?episodeType=first",
	)
}
