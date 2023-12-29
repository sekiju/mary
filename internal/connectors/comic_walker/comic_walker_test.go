package comic_walker

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestComicWalker(t *testing.T) {
	test_utils.ReaderTest(t, New(), "tests/assets/comic_walker.jpg",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189954777350688798/00.jpg",
		"https://comic-walker.com/viewer/?tw=2&dlcl=ja&cid=KDCW_MF10204507010001_68")
}
