package comic_walker

type EpisodeResult struct {
	Episode struct {
		Id       string `json:"id"`
		Code     string `json:"code"`
		Title    string `json:"title"`
		Internal struct {
			EpisodeNo int `json:"episodeNo"`
		} `json:"internal"`
	} `json:"episode"`
}

type ViewerResult struct {
	Manuscripts []struct {
		DrmHash     string `json:"drmHash"`
		DrmImageUrl string `json:"drmImageUrl"`
		Page        int    `json:"page"`
	} `json:"manuscripts"`
}

type ViewerJumpForwardResult struct {
	Episode *struct {
		Id       string `json:"id"`
		Code     string `json:"code"`
		Title    string `json:"title"`
		Internal struct {
			EpisodeNo int `json:"episodeNo"`
		} `json:"internal"`
	} `json:"episode"`
}
