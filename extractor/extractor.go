package extractor

import (
	"fmt"

	"github.com/sekiju/mdl/config"
	"github.com/sekiju/mdl/extractor/registry"
	"github.com/sekiju/mdl/sdk/manga"
)

// getSession matches hostname exactly against config.Params.File.Sites, with
// no subdomain/suffix fallback. A site served from a variable subdomain will
// silently fall through to "unsupported website" in NewExtractor unless this
// is revisited.
func getSession(hostname string) *string {
	if config.Params.Runtime.PrimaryCookie != nil {
		return config.Params.Runtime.PrimaryCookie
	}
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
