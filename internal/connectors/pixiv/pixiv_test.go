package pixiv

import (
	"559/internal/utils/testing_utils"
	"testing"
)

func TestPixiv_Chapter(t *testing.T) {
	testing_utils.Connector(
		t,
		New(),
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189991200262983710/104260548_p0.png",
		"https://www.pixiv.net/en/artworks/104260548",
	)
}
