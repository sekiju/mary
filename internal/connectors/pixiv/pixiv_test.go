package pixiv

import (
	"mary/internal/utils"
	"testing"
)

func TestPixiv_Chapter(t *testing.T) {
	utils.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/pixiv.png",
		"https://www.pixiv.net/en/artworks/104260548",
	)
}
