package manga_bilibili

import (
	"559/internal/static"
	"net/url"
	"testing"
)

func TestMangaBiliBili_ResolveType(t *testing.T) {
	urls := map[static.UrlType]string{
		static.UrlTypeBook:    "https://manga.bilibili.com/detail/mc26505",
		static.UrlTypeChapter: "https://manga.bilibili.com/mc26505/312836?from=manga_detail",
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
