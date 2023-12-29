package pixiv

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestPixiv(t *testing.T) {
	test_utils.ReaderTest(t, New(), "tests/assets/pixiv.png",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189991200262983710/104260548_p0.png",
		"https://www.pixiv.net/en/artworks/104260548")
}
