package registry

import (
	"559/internal/connectors"
	"559/internal/connectors/comic_walker"
	"559/internal/connectors/fod"
	"559/internal/connectors/giga_viewer"
	"559/internal/connectors/newtype"
	"559/internal/connectors/pixiv"
	"559/internal/connectors/speed_binb"
	"fmt"
	"sync"
)

type ReadersRegistry struct {
	readers map[string]connectors.Connector
	mu      sync.RWMutex
}

var Default = &ReadersRegistry{
	readers: make(map[string]connectors.Connector),
}

func (r *ReadersRegistry) Add(n connectors.Connector) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.readers[n.Context().Domain] = n
}

func (r *ReadersRegistry) All() []connectors.Connector {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v := make([]connectors.Connector, len(r.readers))
	for _, value := range r.readers {
		v = append(v, value)
	}

	return v
}

func (r *ReadersRegistry) FindReaderByDomain(domain string) (connectors.Connector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parser, exists := r.readers[domain]
	if !exists {
		return parser, fmt.Errorf("unknown website")
	}

	return parser, nil
}

//type TemplateNewFunc func(domain string) connectors.Connector

//func (r *ReadersRegistry) MassiveAdd(tfunc TemplateNewFunc, domains []string) {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	for _, domain := range domains {
//		reader := tfunc(domain)
//		r.connectors[reader.Domain()] = reader
//	}
//}

type Config map[string]map[string]interface{}

func init() {
	Default.Add(fod.New())
	Default.Add(comic_walker.New())
	Default.Add(pixiv.New())
	Default.Add(newtype.New())

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
		Default.Add(giga_viewer.New(domain))
	}

	speedBinbWebsites := []string{
		"storia.takeshobo.co.jp",
		"www.comic-valkyrie.com",
		"www.cmoa.jp",
		"yanmaga.jp",
	}
	for _, domain := range speedBinbWebsites {
		Default.Add(speed_binb.New(domain))
	}
}
