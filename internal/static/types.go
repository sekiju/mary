package static

import "net/url"

type ImageFn func() ([]byte, error)
type Image struct {
	FileName string
	ImageFn  ImageFn
}

func NewImage(fileName string, fn *ImageFn) Image {
	return Image{
		FileName: fileName,
		ImageFn:  *fn,
	}
}

type AuthorizationStatus string

const (
	AuthorizationStatusRequired     AuthorizationStatus = "REQUIRED"
	AuthorizationStatusOptional                         = "OPTIONAL"
	AuthorizationStatusForBookmarks                     = "FOR_BOOKMARKS"
	AuthorizationStatusNay                              = "NOT_AVAILABLE"
)

type ConnectorData struct {
	Domain               string
	AuthorizationStatus  AuthorizationStatus
	ChapterListAvailable bool
}

type Book struct {
	Title    string    `json:"title"`
	Cover    *string   `json:"cover"`
	Chapters []Chapter `json:"chapters"`
}

type Chapter struct {
	ID    any    `json:"id"`
	Title string `json:"title"`
	Error error  `json:"error"`
}

type UrlType string

const (
	UrlTypeBook    UrlType = "BOOK"
	UrlTypeChapter         = "CHAPTER"
	UrlTypeShared          = "SHARED"
)

type Connector interface {
	Data() *ConnectorData
	ResolveType(uri url.URL) (UrlType, error)
	Book(uri url.URL) (*Book, error)
	Chapter(uri url.URL) (*Chapter, error)
	Pages(chapterID any, imageChan chan<- Image) error
}
