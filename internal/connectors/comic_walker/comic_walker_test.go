package comic_walker

import (
	"559/internal/static"
	"559/internal/utils/testing_utils"
	"net/url"
	"testing"
)

func TestComicWalker_ResolveType(t *testing.T) {
	urls := map[static.UrlType]string{
		static.UrlTypeBook:    "https://comic-walker.com/contents/detail/KDCW_AP01204433010000_68/",
		static.UrlTypeChapter: "https://comic-walker.com/viewer/?tw=2&dlcl=ja&cid=KDCW_AP01204433010001_68",
	}

	connector := New()

	for k, v := range urls {
		uri, _ := url.Parse(v)
		resolveType, err := connector.ResolveType(*uri)
		if err != nil {
			t.Error(err)
		}

		if resolveType != k {
			t.Error("mismatched url type")
		}
	}
}

func TestComicWalker_Chapter(t *testing.T) {
	testing_utils.Connector(
		t,
		New(),
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189954777350688798/00.jpg",
		"https://comic-walker.com/viewer/?tw=2&dlcl=ja&cid=KDCW_MF10204507010001_68",
	)
}
