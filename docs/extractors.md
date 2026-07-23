# Writing a site extractor

This guide walks through building a new extractor for `mdl`. We'll use the
built-in `cmoa` extractor as a reference — it's the simplest of the six
built-in extractors and shows the minimal shape every extractor needs.

## Minimal extractor contract

Every extractor must:

1. **Implement `sdk/manga.Extractor`** — three methods:
   `FindChapters`, `FindChapter`, `FindChapterPages` (all take
   `context.Context` as first argument) plus `SetSettings`.

2. **Embed `sdk/manga/pluginutil.Base`** — stores settings and
   provides `ApplyCookie(req)` for cookie-authenticated requests.

3. **Register via `extractor/registry.Register` from an `init()`
   function** — the standard pattern is
   `registry.Register(hostname, registry.WithSession(New))`.

4. **Have a `New() (manga.Extractor, error)` constructor** —
   `WithSession` calls it to create a fresh instance.

## Worked example: cmoa

```go
package cmoa

import (
 "context"

 "github.com/knst0/mdl/extractor/registry"
 "github.com/knst0/mdl/sdk/manga"
 "github.com/knst0/mdl/sdk/manga/pluginutil"

 "resty.dev/v3"
)

// Extractor embeds pluginutil.Base for settings storage + ApplyCookie.
type Extractor struct {
 pluginutil.Base
}

// FindChapters — cmoa does not support chapter listing.
func (e *Extractor) FindChapters(ctx context.Context, URL string) ([]*manga.Chapter, error) {
 return nil, manga.ErrChapterListingUnsupported
}

// SupportsChapterListing lets callers check via type assertion.
func (e *Extractor) SupportsChapterListing() bool {
 return false
}

// FindChapter parses the chapter ID from the URL and returns a Chapter.
func (e *Extractor) FindChapter(ctx context.Context, URL string) (*manga.Chapter, error) {
 ID, err := extractViewerID(URL)
 if err != nil {
  return nil, err
 }

 return &manga.Chapter{
  ID:  ID,
  URL: URL,
 }, nil
}

// FindChapterPages downloads the chapter pages. Cookie is required.
func (e *Extractor) FindChapterPages(ctx context.Context, chapter *manga.Chapter) ([]*manga.Page, error) {
 if e.Settings.Cookie == nil {
  return nil, manga.ErrCredentialsRequired
 }

 req := resty.New().R().SetContext(ctx)
 e.ApplyCookie(req)
 // ...fetch and build []*manga.Page
}

// init registers this extractor for its hostname(s).
func init() {
 registry.Register("www.cmoa.jp", registry.WithSession(New))
}

// New constructs a new Extractor with default settings.
func New() (manga.Extractor, error) {
 return &Extractor{Base: pluginutil.Base{Settings: &manga.Settings{}}}, nil
}
```

## Key types

### `manga.Extractor`

```go
type Extractor interface {
 FindChapters(ctx context.Context, URL string) ([]*Chapter, error)
 FindChapter(ctx context.Context, URL string) (*Chapter, error)
 FindChapterPages(ctx context.Context, chapter *Chapter) ([]*Page, error)
 SetSettings(settings Settings)
}
```

### `manga.Chapter`

```go
type Chapter struct {
 ID      string // site-specific chapter identifier
 Number  string // human-readable number (e.g. "12")
 Title   string
 Index   uint   // zero-based position in the chapter list
 URL     string
 MangaID string // site-specific manga/series identifier
}
```

### `manga.Page`

```go
type Page struct {
 URL      string            // image URL to download
 Filename string            // output filename (use pluginutil.BuildPages)
 Index    uint              // zero-based page index
 Headers  map[string]string // optional request headers
 Decode   DecodeFunc        // optional post-download decoder (e.g. for DRM)
}
```

### `manga.Settings`

```go
type Settings struct {
 Cookie *string // nil means unauthenticated
}
```

Use `pluginutil.Base.ApplyCookie(req)` to attach the cookie to any
authenticated request. Check `e.Settings.Cookie == nil` and return
`manga.ErrCredentialsRequired` when cookie is required but missing.

## Optional interfaces

### `GenerateCookieFeature`

Implement `GenerateCookie() (string, error)` if the site can generate
guest session cookies automatically. `WithSession` calls this when no
cookie is configured.

### `ChapterListingFeature`

Implement `SupportsChapterListing() bool` (return `true`) if your
extractor can list all chapters of a manga. Callers type-assert before
calling `FindChapters`.

## Page building

Use `pluginutil.BuildPages` for consistent file naming:

```go
return pluginutil.BuildPages(len(rawPages), ".jpg", func(i int, filename string) (*manga.Page, error) {
 return &manga.Page{
  Index:    uint(i),
  URL:      rawPages[i].Src,
  Filename: filename,
 }, nil
})
```

## Registration

```go
import "github.com/knst0/mdl/extractor/registry"

func init() {
 registry.Register("www.example.com", registry.WithSession(New))
}
```

- **`registry.WithSession(New)`** wraps your constructor with cookie
  resolution: config cookie first, then `GenerateCookieFeature` fallback.
- **Duplicate hostname registrations panic** — this is a compile-time
  programming error, not a runtime condition.
- For an extractor that handles multiple hostnames, call `Register`
  once per hostname.

## Wiring it into a binary

If you're building your own `main` that includes custom extractors:

```go
package main

import (
 _ "your.module/path/extractor/mysite" // triggers init() → Register
 "github.com/knst0/mdl/cmd/cli"
)

func main() {
 cli.Execute()
}
```

The built-in extractors are wired the same way — via blank imports in
`extractor/builtin.go`.

## Versioning

`sdk/manga` currently ships inside the main `github.com/knst0/mdl`
module at whatever version the module is tagged. There is no
independent versioning story yet (Changesets-style release tracking is
planned — see Phase 4 item 4.3 of `REFACTOR_PLAN.md`). Until then,
breaking changes to `sdk/manga` types or interfaces will be called out
in commit messages and release notes.

Splitting `sdk/manga` into its own Go module (separate `go.mod`) is
deferred — revisit once automated changelog tooling lands and there is
evidence of real third-party extractor authors.
