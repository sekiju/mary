package extractor

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sekiju/mdl/config"
	"github.com/sekiju/mdl/extractor/cmoa"
	"github.com/sekiju/mdl/extractor/comic_walker"
	"github.com/sekiju/mdl/extractor/corocoro"
	"github.com/sekiju/mdl/extractor/ganma"
	"github.com/sekiju/mdl/extractor/storia_takeshobo"
	"github.com/sekiju/mdl/extractor/template/giga_viewer"
	"github.com/sekiju/mdl/sdk/manga"
)

type Factory func(cookie *string) (manga.Extractor, error)

var domainRegistry = map[string]Factory{
	"comic-walker.com":          withSession(comic_walker.New),
	"shonenjumpplus.com":        withSession(giga_viewer.New),
	"comic-zenon.com":           withSession(giga_viewer.New),
	"pocket.shonenmagazine.com": withSession(giga_viewer.New),
	"comic-gardo.com":           withSession(giga_viewer.New),
	"magcomi.com":               withSession(giga_viewer.New),
	"tonarinoyj.jp":             withSession(giga_viewer.New),
	"comic-ogyaaa.com":          withSession(giga_viewer.New),
	"comic-action.com":          withSession(giga_viewer.New),
	"comic-days.com":            withSession(giga_viewer.New),
	"comic-growl.com":           withSession(giga_viewer.New),
	"comic-earthstar.com":       withSession(giga_viewer.New),
	"comicborder.com":           withSession(giga_viewer.New),
	"comic-trail.com":           withSession(giga_viewer.New),
	"kuragebunch.com":           withSession(giga_viewer.New),
	"viewer.heros-web.com":      withSession(giga_viewer.New),
	"www.sunday-webry.com":      withSession(giga_viewer.New),
	"www.cmoa.jp":               withSession(cmoa.New),
	"www.corocoro.jp":           withSession(corocoro.New),
	"storia.takeshobo.co.jp":    withSession(storia_takeshobo.New),
	"ganma.jp":                  withSession(ganma.New),
}

func withSession(fn func() (manga.Extractor, error)) Factory {
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
	factory, exists := domainRegistry[hostname]
	if !exists {
		return nil, fmt.Errorf("unsupported website: %s", hostname)
	}
	return factory(getSession(hostname))
}
