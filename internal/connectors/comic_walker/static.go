package comic_walker

type BookResponse struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Result struct {
			Id   int `json:"id"`
			Meta struct {
				Title             string        `json:"title"`
				DisplayAuthorName string        `json:"display_author_name"`
				Authors           []interface{} `json:"authors"`
				Description       string        `json:"description"`
				PlayerType        int           `json:"player_type"`
				ContentType       string        `json:"content_type"`
				Official          interface{}   `json:"official"`
				SerialStatus      int           `json:"serial_status"`
				Rating            struct {
					Violence interface{} `json:"violence"`
					Adult    interface{} `json:"adult"`
				} `json:"rating"`
				Expressions struct {
					BoysLove bool `json:"boys_love"`
				} `json:"expressions"`
				Category string `json:"category"`
				Counter  struct {
					View     int `json:"view"`
					Comment  int `json:"comment"`
					Favorite int `json:"favorite"`
					Episode  int `json:"episode"`
				} `json:"counter"`
				PagePadding       string      `json:"page_padding"`
				CreatedAt         int         `json:"created_at"`
				UpdatedAt         int         `json:"updated_at"`
				UpdateScheduledAt interface{} `json:"update_scheduled_at"`
				IconUrl           string      `json:"icon_url"`
				ThumbnailUrl      string      `json:"thumbnail_url"`
				MainImageUrl      string      `json:"main_image_url"`
				ShareUrl          string      `json:"share_url"`
			} `json:"meta"`
		} `json:"result"`
		Extra struct {
			Direction string `json:"direction"`
		} `json:"extra"`
	} `json:"data"`
}

type ChapterResponse struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Result struct {
			Title string `json:"title"`
		} `json:"result"`
		Extra struct {
			Content struct {
				Title    string `json:"title"`
				ShareUrl string `json:"share_url"`
			} `json:"content"`
			Direction string `json:"direction"`
		} `json:"extra"`
	} `json:"data"`
}

type FramesResponse struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data struct {
		Result []struct {
			Id   int `json:"id"`
			Meta struct {
				Width     int     `json:"width"`
				Height    int     `json:"height"`
				SourceUrl string  `json:"source_url"`
				DrmHash   *string `json:"drm_hash"`
				Duration  int     `json:"duration"`
				LinkUrl   *string `json:"link_url"`
				Resource  struct {
					Bgm interface{} `json:"bgm"`
					Se  interface{} `json:"se"`
				} `json:"resource"`
				IsSpread bool `json:"is_spread"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"data"`
}
