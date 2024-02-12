package comic_webnewtype

import (
	"559/internal/static"
	"559/internal/utils/testing_utils"
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
	testing_utils.Connector(
		t,
		New(),
		"https://cdn.discordapp.com/attachments/1157977166022189086/1190705156078194740/mankitsu_12_01-a989e6fc-4285-4a38-a2cb-2110f40d09d6.jpg",
		"https://comic.webnewtype.com/contents/mankitsu/120/",
	)
}
