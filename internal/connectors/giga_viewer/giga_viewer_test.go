package giga_viewer

import (
	"559/internal/utils/test_utils"
	"testing"
)

func TestShonenJumpPlus(t *testing.T) {
	test_utils.ReaderTest(t, New("shonenjumpplus.com"), "tests/assets/shonenjumpplus.jpg",
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189993549362630727/02.jpg",
		"https://shonenjumpplus.com/episode/10834108156648240735")
}
