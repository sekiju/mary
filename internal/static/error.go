package static

import "errors"

var (
	UnsupportedWebsiteError = errors.New("unsupported website")
	LoginRequiredError      = errors.New("login required")
	UnknownUrlType          = errors.New("unknown URL type")
)
