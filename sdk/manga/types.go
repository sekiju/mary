package manga

type (
	Extractor interface {
		FindChapters(URL string) ([]*Chapter, error)
		FindChapter(URL string) (*Chapter, error)
		FindChapterPages(chapter *Chapter) ([]*Page, error)
		SetSettings(settings Settings)
	}

	GenerateCookieFeature interface {
		GenerateCookie() (string, error)
	}

	// ChapterListingFeature is implemented by extractors that can list all
	// chapters of a manga via FindChapters. Extractors that only support
	// resolving a single chapter URL (FindChapter) should not implement it,
	// so callers can check via a type assertion before invoking FindChapters.
	ChapterListingFeature interface {
		SupportsChapterListing() bool
	}

	Settings struct {
		Cookie *string
	}

	Chapter struct {
		ID      string `json:"id"`
		Number  string `json:"number"`
		Title   string `json:"title"`
		Index   uint   `json:"index"`
		URL     string `json:"url"`
		MangaID string `json:"mangaId"`
	}

	DecodeFunc func(b []byte) ([]byte, error)

	Page struct {
		URL      string            `json:"url"`
		Filename string            `json:"filename"`
		Index    uint              `json:"index"`
		Headers  map[string]string `json:"headers"`
		Decode   DecodeFunc
	}
)
