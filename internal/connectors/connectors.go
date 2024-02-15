package connectors

import (
	"559/internal/connectors/comic_walker"
	"559/internal/connectors/comic_webnewtype"
	"559/internal/connectors/fod"
	"559/internal/connectors/giga_viewer"
	"559/internal/connectors/manga_bilibili"
	"559/internal/connectors/pixiv"
	"559/internal/connectors/speed_binb/cmoa"
	"559/internal/connectors/speed_binb/comic_meteor"
	"559/internal/connectors/speed_binb/comic_valkyrie"
	"559/internal/connectors/speed_binb/storia_takeshobo"
	"559/internal/connectors/speed_binb/yanmaga"
	"559/internal/static"
)

var state = make(map[string]static.Connector)

func FindByDomain(domain string) (static.Connector, error) {
	connector, ok := state[domain]
	if !ok {
		return nil, static.UnsupportedWebsiteErr
	}

	return connector, nil
}

func Add(connectors ...static.Connector) {
	for _, connector := range connectors {
		state[connector.Data().Domain] = connector
	}
}

func init() {
	Add(
		comic_walker.New(),
		pixiv.New(),
		comic_webnewtype.New(),
		manga_bilibili.New(),
		fod.New(),
		comic_valkyrie.New(),
		cmoa.New(),
		storia_takeshobo.New(),
		yanmaga.New(),
		comic_meteor.New(),
	)

	gigaViewerWebsites := []string{
		"shonenjumpplus.com",
		"pocket.shonenmagazine.com",
		"comic-action.com",
		"comic-days.com",
		"comic-growl.com",
		"comic-earthstar.com",
		"comic-gardo.com",
		"comic-trail.com",
		"comic-zenon.com",
		"comicborder.com",
		"kuragebunch.com",
		"magcomi.com",
		"tonarinoyj.jp",
		"viewer.heros-web.com",
		"www.sunday-webry.com",
		"comicbushi-web.com",
	}
	for _, domain := range gigaViewerWebsites {
		Add(giga_viewer.New(domain))
	}
}
