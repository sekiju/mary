package registry

import (
	"559/internal/readers"
	"559/internal/readers/comic_walker"
	"559/internal/readers/fod"
	"559/internal/readers/giga_viewer"
	"559/internal/readers/pixiv"
	"559/internal/utils"
	"encoding/json"
	"fmt"
	"sync"
)

type ReadersRegistry struct {
	readers map[string]readers.Reader
	mu      sync.RWMutex
}

var Default = &ReadersRegistry{
	readers: make(map[string]readers.Reader),
}

func (r *ReadersRegistry) Add(n readers.Reader) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.readers[n.Context().Domain] = n
}

func (r *ReadersRegistry) All() []readers.Reader {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v := make([]readers.Reader, len(r.readers))
	for _, value := range r.readers {
		v = append(v, value)
	}

	return v
}

func (r *ReadersRegistry) FindReaderByDomain(domain string) (readers.Reader, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parser, exists := r.readers[domain]
	if !exists {
		return parser, fmt.Errorf("unknown website")
	}

	return parser, nil
}

//type TemplateNewFunc func(domain string) readers.Reader

//func (r *ReadersRegistry) MassiveAdd(tfunc TemplateNewFunc, domains []string) {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	for _, domain := range domains {
//		reader := tfunc(domain)
//		r.readers[reader.Domain()] = reader
//	}
//}

type Config map[string]map[string]interface{}

func init() {
	Default.Add(fod.New())
	Default.Add(comic_walker.New())
	Default.Add(pixiv.New())

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
	}
	for _, domain := range gigaViewerWebsites {
		Default.Add(giga_viewer.New(domain))
	}

	configFile, err := utils.ReadFile("settings.json")
	if err == nil {
		var config Config
		err := json.Unmarshal(configFile, &config)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		for domain, cfg := range config {
			r, err := Default.FindReaderByDomain(domain)
			if err != nil {
				fmt.Printf("Unknown website %q in settings.json\n", domain)
				continue
			}

			for k, v := range cfg {
				r.Context().Data[k] = v
			}
		}
	} else {
		fmt.Println("Create settings.json file for download private books")
	}
}
