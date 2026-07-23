package extractor

import (
	"fmt"

	"github.com/knst0/mdl/config"
	"github.com/knst0/mdl/extractor/registry"
	"github.com/knst0/mdl/sdk/manga"
)

// getSession matches hostname exactly against config.Params.File.Sites, with
// no subdomain/suffix fallback. A site served from a variable subdomain will
// silently fall through to "unsupported website" in NewExtractor unless this
// is revisited.
func getSession(hostname string) *string {
	if site, exists := config.Params.File.Sites[hostname]; exists && site.Cookie != nil {
		return site.Cookie
	}
	return nil
}

func NewExtractor(hostname string) (manga.Extractor, error) {
	factory, exists := registry.Lookup(hostname)
	if !exists {
		return nil, fmt.Errorf("unsupported website: %s", hostname)
	}
	return factory(getSession(hostname))
}
