package comic_walker

type FramesResponse struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Result []struct {
			Id   int `json:"id"`
			Meta struct {
				Width     int         `json:"width"`
				Height    int         `json:"height"`
				SourceUrl string      `json:"source_url"`
				DrmHash   string      `json:"drm_hash"`
				Duration  int         `json:"duration"`
				LinkUrl   interface{} `json:"link_url"`
				Resource  struct {
					Bgm interface{} `json:"bgm"`
					Se  interface{} `json:"se"`
				} `json:"resource"`
				IsSpread bool `json:"is_spread"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"data"`
}
