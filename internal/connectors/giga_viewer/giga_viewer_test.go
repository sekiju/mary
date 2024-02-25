package giga_viewer

import (
	"mary/tools"
	"testing"
)

func TestGigaViewer_Chapter_ShonenJumpPlus(t *testing.T) {
	tools.TestConnector(
		t,
		New("shonenjumpplus.com"),
		"https://host.kireyev.org/mary-files/shonenjumpplus.jpg?",
		"https://shonenjumpplus.com/episode/10834108156648240735",
	)
}
