package pixiv

type IllustResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    struct {
		Id         string `json:"id"`
		Title      string `json:"title"`
		IllustType int    `json:"illustType"`
		Urls       struct {
			Mini     *string `json:"mini,omitempty"`
			Thumb    *string `json:"thumb,omitempty"`
			Small    *string `json:"small,omitempty"`
			Regular  *string `json:"regular,omitempty"`
			Original *string `json:"original,omitempty"`
		} `json:"urls"`
		Tags struct {
			AuthorId string `json:"authorId"`
			IsLocked bool   `json:"isLocked"`
			Tags     []struct {
				Tag         string `json:"tag"`
				Locked      bool   `json:"locked"`
				Deletable   bool   `json:"deletable"`
				UserId      string `json:"userId,omitempty"`
				UserName    string `json:"userName,omitempty"`
				Romaji      string `json:"romaji,omitempty"`
				Translation struct {
					En string `json:"en"`
				} `json:"translation,omitempty"`
			} `json:"tags"`
			Writable bool `json:"writable"`
		} `json:"tags"`
		PageCount   int       `json:"pageCount"`
		NoLoginData *struct{} `json:"noLoginData,omitempty"`
	} `json:"body"`
}

type PagesResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    []struct {
		Urls struct {
			ThumbMini string `json:"thumb_mini"`
			Small     string `json:"small"`
			Regular   string `json:"regular"`
			Original  string `json:"original"`
		} `json:"urls"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"body"`
}
