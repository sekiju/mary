package registry

import (
	"559/internal/readers"
	"559/internal/readers/comic_walker"
	"559/internal/readers/fod"
	"559/internal/readers/shonenjumpplus"
	"sync"
)

// todo: exception for unsupported website

type ReadersRegistry struct {
	parsers map[string]readers.Parser
	mu      sync.RWMutex
}

var Default = &ReadersRegistry{
	parsers: make(map[string]readers.Parser),
}

func (r *ReadersRegistry) Add(parser readers.Parser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[parser.Details().ID] = parser
}

func (r *ReadersRegistry) All() []readers.Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v := make([]readers.Parser, 0, len(r.parsers))

	for _, value := range r.parsers {
		v = append(v, value)
	}

	return v
}

func (r *ReadersRegistry) FindParserByID(id string) readers.Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.parsers[id]
}

func (r *ReadersRegistry) FindParserByDomain(domain string) readers.Parser {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := range r.parsers {
		if r.parsers[i].Details().Domain == domain {
			return r.parsers[i]
		}
	}

	return nil
}

// todo: move configs to main file

func init() {
	Default.Add(&fod.Fod{DebugKeys: true, Session: "VWdx8id9R0XHjVpvs7s754CxGJpBBl9ZCHCqL1yF"})
	Default.Add(&shonenjumpplus.ShonenJumpPlus{})
	Default.Add(&comic_walker.ComicWalker{})
}
