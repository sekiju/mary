package pixiv

import (
	"mary/test"
	"testing"
)

func TestPixiv_Chapter(t *testing.T) {
	test.Connector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/pixiv.png",
		"https://www.pixiv.net/en/artworks/104260548",
	)
}
