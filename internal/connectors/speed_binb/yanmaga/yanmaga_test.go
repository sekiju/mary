package yanmaga

import (
	"559/internal/utils/testing_utils"
	"testing"
)

func TestYanmaga_Chapter(t *testing.T) {
	testing_utils.Connector(
		t,
		New(),
		"https://host.kireyev.org/559/yanmaga.jpg",
		"https://yanmaga.jp/viewer/comics/その無能、実は世界最強の魔法使い_無能と蔑まれ、貴族家から追い出されたが、ギフト転生者が覚醒して前世の能力が蘇った/546838aaf7e2832b48c199c7af660e6d?cid=06A0000000000425482T",
	)
}
