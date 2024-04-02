package comic_walker

import (
	"mary/test"
	"testing"
)

func TestComicWalker_Chapter(t *testing.T) {
	test.Connector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/comic_walker.webp",
		"https://comic-walker.com/detail/KC_004526_S/episodes/KC_0045260000100011_E?episodeType=first",
	)
}
