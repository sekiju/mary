package giga_viewer

import (
	"mary/internal/utils"
	"testing"
)

func TestGigaViewer_Chapter_ShonenJumpPlus(t *testing.T) {
	utils.TestConnector(
		t,
		New("shonenjumpplus.com"),
		"https://host.kireyev.org/mary-files/shonenjumpplus.jpg?",
		"https://shonenjumpplus.com/episode/10834108156648240735",
	)
}
