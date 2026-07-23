package manga

import "errors"

var (
	ErrMethodUnimplemented       = errors.New("method unimplemented")
	ErrMangaNotFound             = errors.New("manga not found")
	ErrChapterNotFound           = errors.New("chapter not found")
	ErrPaidChapter               = errors.New("chapter is paid")
	ErrInvalidURLFormat          = errors.New("invalid URL format")
	ErrInvalidChapterURL         = errors.New("cannot extract chapterID from this url")
	ErrCredentialsRequired       = errors.New("credentials required: set site.'domain'.cookie field in config.json")
	ErrChapterListingUnsupported = errors.New("this extractor does not support listing chapters")
	ErrMalformedChapterData      = errors.New("chapter data from source is malformed or corrupted")
	ErrUnsupportedContentFormat  = errors.New("content format returned by source is not supported")
)
