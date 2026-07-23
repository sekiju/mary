// Package registry provides a plugin registry for manga site extractors.
// Extractors call Register from their init() functions; callers resolve
// extractors via NewExtractor.
package registry

import (
	"fmt"
	"sort"

	"github.com/knst0/mdl/sdk/manga"
	"github.com/rs/zerolog/log"
)

type Factory func(cookie *string) (manga.Extractor, error)

var domainRegistry = map[string]Factory{}

// Register adds a factory for hostname. A duplicate registration panics
// because it is a compile-time programming error, not a runtime condition.
func Register(hostname string, factory Factory) {
	if _, exists := domainRegistry[hostname]; exists {
		panic(fmt.Sprintf("extractor/registry: duplicate registration for hostname %q", hostname))
	}
	domainRegistry[hostname] = factory
}

// Lookup returns the factory registered for hostname, or false if none.
func Lookup(hostname string) (Factory, bool) {
	f, ok := domainRegistry[hostname]
	return f, ok
}

// Hostnames returns all registered hostnames, sorted.
func Hostnames() []string {
	hosts := make([]string, 0, len(domainRegistry))
	for h := range domainRegistry {
		hosts = append(hosts, h)
	}
	sort.Strings(hosts)
	return hosts
}

func WithSession(fn func() (manga.Extractor, error)) Factory {
	return func(cookie *string) (manga.Extractor, error) {
		ext, err := fn()
		if err != nil {
			return nil, err
		}

		if cookie != nil {
			ext.SetSettings(manga.Settings{Cookie: cookie})
		} else {
			cookieGenerator, ok := ext.(manga.GenerateCookieFeature)
			if ok {
				generatedCookie, err := cookieGenerator.GenerateCookie()
				if err != nil {
					return nil, err
				}

				log.Info().Msgf("Cookie generated >>> %s", generatedCookie)

				ext.SetSettings(manga.Settings{Cookie: &generatedCookie})
			}
		}

		return ext, nil
	}
}
