package shonenmagazine

type apiResponse struct {
	Status       string `json:"status"`
	ResponseCode int    `json:"response_code"`
	ErrorMessage string `json:"error_message"`
}

type viewerJSON struct {
	apiResponse
	TitleID      int64    `json:"title_id"`
	ScrambleSeed string   `json:"scramble_seed"`
	PageList     []string `json:"page_list"`
}

type episodeEntry struct {
	EpisodeID   int64  `json:"episode_id"`
	EpisodeName string `json:"episode_name"`
	Index       int    `json:"index"`
	TitleID     int64  `json:"title_id"`
}

type episodeListJSON struct {
	apiResponse
	EpisodeList []episodeEntry `json:"episode_list"`
}

type titleDetailJSON struct {
	apiResponse
	WebTitle struct {
		EpisodeIDList []int64 `json:"episode_id_list"`
	} `json:"web_title"`
}
