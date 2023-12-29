package fod

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestFod(t *testing.T) {
	test_utils.ReaderTest(t, New(), "tests/assets/fod.jpg",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189988876429820065/01.jpg",
		"https://manga.fod.fujitv.co.jp/books/1094816/BT000109481600100101")
}
