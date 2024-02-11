package connectors

import (
	"559/internal/connectors/comic_walker"
	"559/internal/connectors/giga_viewer"
	"559/internal/connectors/pixiv"
	"559/internal/static"
)

var state = make(map[string]static.Connector)

func FindByDomain(domain string) (static.Connector, error) {
	connector, ok := state[domain]
	if !ok {
		return nil, static.UnsupportedWebsiteError
	}

	return connector, nil
}

func Add(connectors ...static.Connector) {
	for _, connector := range connectors {
		state[connector.Data().Domain] = connector
	}
}

func init() {
	Add(comic_walker.New(), pixiv.New())

	//Add(fod.New(), comic_walker.New(), pixiv.New(), newtype.New())

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
	//
	//speedBinbWebsites := []string{
	//	"storia.takeshobo.co.jp",
	//	"www.comic-valkyrie.com",
	//	"www.cmoa.jp",
	//	"yanmaga.jp",
	//}
	//for _, domain := range speedBinbWebsites {
	//	Add(speed_binb.New(domain))
	//}
}
