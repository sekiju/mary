package pixiv

import (
	"mary/tools"
	"testing"
)

func TestPixiv_Chapter(t *testing.T) {
	tools.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/pixiv.png",
		"https://www.pixiv.net/en/artworks/104260548",
	)
}
