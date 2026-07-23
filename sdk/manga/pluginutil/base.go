package pluginutil

import (
	"github.com/knst0/mdl/sdk/manga"
	"resty.dev/v3"
)

// Base provides the settings storage and SetSettings implementation shared
// by every extractor. Extractors embed Base instead of redeclaring the
// settings field and method.
type Base struct {
	Settings *manga.Settings
}

func (b *Base) SetSettings(settings manga.Settings) {
	b.Settings = &settings
}

// ApplyCookie sets the Cookie header on req from b.Settings.Cookie if set.
func (b *Base) ApplyCookie(req *resty.Request) {
	if b.Settings.Cookie != nil {
		req.SetHeader("Cookie", *b.Settings.Cookie)
	}
}
