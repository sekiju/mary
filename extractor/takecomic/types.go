package takecomic

// EpisodeResponse is the JSON returned by GET /api/episodes/{id}.
type EpisodeResponse struct {
	Episode Episode `json:"episode"`
}

type Episode struct {
	ContentID int     `json:"contentId"`
	ID        string  `json:"id"`
	IndexID   int     `json:"indexId"`
	Series    Series  `json:"series"`
	Summary   Summary `json:"summary"`
	Content   []struct {
		Type     string `json:"type"`
		ViewerID string `json:"viewerId"`
	} `json:"content"`
	FirstEpisodeID      string `json:"firstEpisodeId"`
	PreviousEpisodeID   string `json:"previousEpisodeId"`
	NextEpisodeID       string `json:"nextEpisodeId"`
	SeriesEpisodeNumber int    `json:"seriesEpisodeNumber"`
}

type Series struct {
	ID      string `json:"id"`
	IndexID int    `json:"indexId"`
	Name    string `json:"name"`
	Status  string `json:"status"`
}

type Summary struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	DatePublished int64  `json:"datePublished"`
	NumLikes      int    `json:"numLikes"`
}

// ContentsInfoResponse is the JSON returned by GET /api/book/contentsInfo.
type ContentsInfoResponse struct {
	TotalPages      int             `json:"totalPages"`
	ScrollDirection string          `json:"scrollDirection"`
	Result          []PageImageInfo `json:"result"`
}

type PageImageInfo struct {
	ImageURL  string `json:"imageUrl"`
	Scramble  string `json:"scramble"` // JSON array like "[6,3,0,14,...]"
	Sort      int    `json:"sort"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	ExpiresOn int64  `json:"expiresOn"`
}
