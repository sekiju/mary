package fod

import (
	"mary/internal/utils"
	"testing"
)

func TestComicNewtype_Chapter(t *testing.T) {
	utils.TestConnector(
		t,
		New(),
		"https://host.kireyev.org/mary-files/manga_fod.jpg",
		"https://manga.fod.fujitv.co.jp/books/1094816/BT000109481600100101",
	)
}
