package giga_viewer

import (
	"mary/test"
	"testing"
)

func TestGigaViewer_Chapter_ShonenJumpPlus(t *testing.T) {
	test.Connector(
		t,
		New("shonenjumpplus.com"),
		"https://host.kireyev.org/mary-files/shonenjumpplus.jpg?",
		"https://shonenjumpplus.com/episode/10834108156648240735",
	)
}
