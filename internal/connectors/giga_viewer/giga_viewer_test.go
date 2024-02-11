package giga_viewer

import (
	"559/internal/utils/testing_utils"
	"testing"
)

func TestGigaViewer_Chapter_ShonenJumpPlus(t *testing.T) {
	testing_utils.Connector(
		t,
		New("shonenjumpplus.com"),
		"https://cdn.discordapp.com/attachments/1157977166022189086/1189993549362630727/02.jpg",
		"https://shonenjumpplus.com/episode/10834108156648240735",
	)
}
