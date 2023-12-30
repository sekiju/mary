package speed_binb

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestPtimg(t *testing.T) {
	test_utils.ReaderTest(t, New("storia.takeshobo.co.jp"), "tests/assets/speed_binb_takeshobo.jpg",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1190016871341232148/00.jpg",
		"https://storia.takeshobo.co.jp/_files/antheia/11_1/")
}
