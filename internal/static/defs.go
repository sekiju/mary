package static

import "errors"

var (
	UnsupportedWebsiteErr           = errors.New("unsupported website")
	LoginRequiredErr                = errors.New("login required")
	UnknownUrlTypeErr               = errors.New("unknown URL type")
	NotFoundErr                     = errors.New("book/chapter not found")
	PaidChapterErr                  = errors.New("paid chapter")
	MassiveDownloaderUnsupportedErr = errors.New("massive downloader unsupported")
)
