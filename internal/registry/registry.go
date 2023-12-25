package registry

import (
	"559/internal/readers"
	"559/internal/readers/comic_walker"
	"559/internal/readers/fod"
	"559/internal/readers/pixiv"
	"559/internal/readers/shonenjumpplus"
	"559/internal/utils"
	"encoding/json"
	"fmt"
	"sync"
)

// todo: exception for unsupported website

type ReadersRegistry struct {
	parsers map[string]readers.Reader
	mu      sync.RWMutex
}

var Default = &ReadersRegistry{
	parsers: make(map[string]readers.Reader),
}

func (r *ReadersRegistry) Add(parser readers.Reader) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[parser.Context().Domain] = parser
}

func (r *ReadersRegistry) All() []readers.Reader {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v := make([]readers.Reader, 0, len(r.parsers))

	for _, value := range r.parsers {
		v = append(v, value)
	}

	return v
}

func (r *ReadersRegistry) FindParserByDomain(domain string) (readers.Reader, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	parser, exists := r.parsers[domain]
	if !exists {
		return parser, fmt.Errorf("unknown website")
	}

	return parser, nil
}

func init() {
	Default.Add(fod.New())
	Default.Add(shonenjumpplus.New())
	Default.Add(comic_walker.New())
	Default.Add(pixiv.New())

	configFile, err := utils.ReadFile("settings.json")
	if err == nil {
		var config map[string]map[string]any
		err := json.Unmarshal(configFile, &config)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		for parserID, cfg := range config {
			parser, err := Default.FindParserByDomain(parserID)
			if err != nil {
				fmt.Printf("Unknown website %q in settings.json\n", parserID)
				continue
			}

			for k, v := range cfg {
				parser.UpdateData(k, v)
			}
		}
	} else {
		fmt.Println("Create settings.json file for download private books")
	}
}
