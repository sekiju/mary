package comic_webnewtype

import (
	"mary/internal/static"
	"mary/internal/utils"
	"net/url"
	"testing"
)

func TestComicNewtype_ResolveType(t *testing.T) {
	urls := map[static.UrlType]string{
		static.UrlTypeBook:    "https://comic.webnewtype.com/contents/mankitsu/",
		static.UrlTypeChapter: "https://comic.webnewtype.com/contents/mankitsu/180/",
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

func TestComicNewtype_Chapter(t *testing.T) {
	utils.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/comic_webnewtype.jpg",
		"https://comic.webnewtype.com/contents/mankitsu/120/",
	)
}
