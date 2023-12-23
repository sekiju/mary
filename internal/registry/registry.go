package registry

import (
	"559/internal/readers"
	"559/internal/readers/comic_walker"
	"559/internal/readers/fod"
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
	r.parsers[parser.Details().ID] = parser
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

func (r *ReadersRegistry) FindParserByID(id string) (readers.Reader, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	parser, exists := r.parsers[id]
	if !exists {
		return parser, fmt.Errorf("unknown parser")
	}

	return parser, nil
}

func (r *ReadersRegistry) FindParserByDomain(domain string) readers.Reader {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.parsers {
		if r.parsers[i].Details().Domain == domain {
			return r.parsers[i]
		}
	}

	return nil
}

type ParserConfig struct {
	Session string `json:"session"`
}

func init() {
	Default.Add(fod.New())
	Default.Add(shonenjumpplus.New())
	Default.Add(comic_walker.New())

	configFile, err := utils.ReadFile("settings.json")
	if err == nil {
		var config map[string]ParserConfig
		err := json.Unmarshal(configFile, &config)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		for parserID, config := range config {
			parser, err := Default.FindParserByID(parserID)
			if err != nil {
				fmt.Println("Invalid parser in settings.json:", err)
			}

			if len(config.Session) > 0 {
				parser.SetSession(config.Session)
			}
		}
	} else {
		fmt.Println("Create settings.json file for download private books")
	}
}
