package manga_bilibili

type GetEpisodeRequest struct {
	ID int `json:"id"`
}

type GetEpisodeResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Title      string `json:"title"`
		ComicId    int    `json:"comic_id"`
		ShortTitle string `json:"short_title"`
		ComicTitle string `json:"comic_title"`
	} `json:"data"`
}

type GetImageIndexRequest struct {
	EpId int `json:"ep_id"`
}

type GetImageIndexResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Path   string `json:"path"`
		Images []struct {
			Path string `json:"path"`
		} `json:"images"`
		Host string `json:"host"`
	} `json:"data"`
}

type ImageTokenRequest struct {
	Urls string `json:"urls"`
}

type ImageTokenResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	} `json:"data"`
}

type ImageTokenError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Meta struct {
		Argument string `json:"argument"`
	} `json:"meta"`
}
