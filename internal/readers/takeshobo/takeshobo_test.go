package takeshobo

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestTakeshobo(t *testing.T) {
	test_utils.ReaderTest(t, New(), "tests/assets/takeshobo.jpg",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1190016871341232148/00.jpg",
		"https://storia.takeshobo.co.jp/_files/antheia/11_1/")
}
