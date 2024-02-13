package fod

import (
	"559/internal/utils/testing_utils"
	"testing"
)

func TestComicNewtype_Chapter(t *testing.T) {
	testing_utils.Connector(
		t,
		New(),
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189988876429820065/01.jpg",
		"https://manga.fod.fujitv.co.jp/books/1094816/BT000109481600100101",
	)
}
