package comic_walker

import (
	"mary/tools"
	"testing"
)

func TestComicWalker_Chapter(t *testing.T) {
	tools.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/comic_walker.webp",
		"https://comic-walker.com/detail/KC_004526_S/episodes/KC_0045260000100011_E?episodeType=first",
	)
}
