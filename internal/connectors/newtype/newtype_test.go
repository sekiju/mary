package newtype

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestNewtype(t *testing.T) {
	test_utils.ReaderTest(t, New(), "tests/assets/newtype.jpg",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1190705156078194740/mankitsu_12_01-a989e6fc-4285-4a38-a2cb-2110f40d09d6.jpg",
		"https://comic.webnewtype.com/contents/mankitsu/120/")
}
