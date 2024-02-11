package pixiv

import "time"

type IllustResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    struct {
		IllustId      string    `json:"illustId"`
		IllustTitle   string    `json:"illustTitle"`
		IllustComment string    `json:"illustComment"`
		Id            string    `json:"id"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		IllustType    int       `json:"illustType"`
		CreateDate    time.Time `json:"createDate"`
		UploadDate    time.Time `json:"uploadDate"`
		Restrict      int       `json:"restrict"`
		XRestrict     int       `json:"xRestrict"`
		Sl            int       `json:"sl"`
		Urls          struct {
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
		Alt         string `json:"alt"`
		UserId      string `json:"userId"`
		UserName    string `json:"userName"`
		UserAccount string `json:"userAccount"`
		UserIllusts map[string]struct {
			Id                      string      `json:"id"`
			Title                   string      `json:"title"`
			IllustType              int         `json:"illustType"`
			XRestrict               int         `json:"xRestrict"`
			Restrict                int         `json:"restrict"`
			Sl                      int         `json:"sl"`
			Url                     string      `json:"url"`
			Description             string      `json:"description"`
			Tags                    []string    `json:"tags"`
			UserId                  string      `json:"userId"`
			UserName                string      `json:"userName"`
			Width                   int         `json:"width"`
			Height                  int         `json:"height"`
			PageCount               int         `json:"pageCount"`
			IsBookmarkable          bool        `json:"isBookmarkable"`
			BookmarkData            interface{} `json:"bookmarkData"`
			Alt                     string      `json:"alt"`
			TitleCaptionTranslation struct {
				WorkTitle   interface{} `json:"workTitle"`
				WorkCaption interface{} `json:"workCaption"`
			} `json:"titleCaptionTranslation"`
			CreateDate time.Time `json:"createDate"`
			UpdateDate time.Time `json:"updateDate"`
			IsUnlisted bool      `json:"isUnlisted"`
			IsMasked   bool      `json:"isMasked"`
			AiType     int       `json:"aiType"`
		} `json:"userIllusts"`
		LikeData             bool          `json:"likeData"`
		Width                int           `json:"width"`
		Height               int           `json:"height"`
		PageCount            int           `json:"pageCount"`
		BookmarkCount        int           `json:"bookmarkCount"`
		LikeCount            int           `json:"likeCount"`
		CommentCount         int           `json:"commentCount"`
		ResponseCount        int           `json:"responseCount"`
		ViewCount            int           `json:"viewCount"`
		BookStyle            any           `json:"bookStyle"`
		IsHowto              bool          `json:"isHowto"`
		IsOriginal           bool          `json:"isOriginal"`
		ImageResponseOutData []interface{} `json:"imageResponseOutData"`
		ImageResponseData    []interface{} `json:"imageResponseData"`
		ImageResponseCount   int           `json:"imageResponseCount"`
		PollData             interface{}   `json:"pollData"`
		SeriesNavData        interface{}   `json:"seriesNavData"`
		DescriptionBoothId   interface{}   `json:"descriptionBoothId"`
		DescriptionYoutubeId interface{}   `json:"descriptionYoutubeId"`
		ComicPromotion       interface{}   `json:"comicPromotion"`
		FanboxPromotion      interface{}   `json:"fanboxPromotion"`
		ContestBanners       []interface{} `json:"contestBanners"`
		IsBookmarkable       bool          `json:"isBookmarkable"`
		BookmarkData         interface{}   `json:"bookmarkData"`
		ContestData          interface{}   `json:"contestData"`
		ZoneConfig           struct {
			Responsive struct {
				Url string `json:"url"`
			} `json:"responsive"`
			Rectangle struct {
				Url string `json:"url"`
			} `json:"rectangle"`
			X500 struct {
				Url string `json:"url"`
			} `json:"500x500"`
			Header struct {
				Url string `json:"url"`
			} `json:"header"`
			Footer struct {
				Url string `json:"url"`
			} `json:"footer"`
			ExpandedFooter struct {
				Url string `json:"url"`
			} `json:"expandedFooter"`
			Logo struct {
				Url string `json:"url"`
			} `json:"logo"`
			Relatedworks struct {
				Url string `json:"url"`
			} `json:"relatedworks"`
		} `json:"zoneConfig"`
		ExtraData struct {
			Meta struct {
				Title              string      `json:"title"`
				Description        string      `json:"description"`
				Canonical          string      `json:"canonical"`
				AlternateLanguages interface{} `json:"alternateLanguages"`
				DescriptionHeader  string      `json:"descriptionHeader"`
				Ogp                struct {
					Description string `json:"description"`
					Image       string `json:"image"`
					Title       string `json:"title"`
					Type        string `json:"type"`
				} `json:"ogp"`
				Twitter struct {
					Description string `json:"description"`
					Image       string `json:"image"`
					Title       string `json:"title"`
					Card        string `json:"card"`
				} `json:"twitter"`
			} `json:"meta"`
		} `json:"extraData"`
		TitleCaptionTranslation struct {
			WorkTitle   interface{} `json:"workTitle"`
			WorkCaption interface{} `json:"workCaption"`
		} `json:"titleCaptionTranslation"`
		IsUnlisted   bool        `json:"isUnlisted"`
		Request      interface{} `json:"request"`
		CommentOff   int         `json:"commentOff"`
		AiType       int         `json:"aiType"`
		ReuploadDate interface{} `json:"reuploadDate"`
		LocationMask bool        `json:"locationMask"`
		NoLoginData  *struct {
			Breadcrumbs struct {
				Successor []interface{} `json:"successor"`
				Current   struct {
					En string `json:"en"`
				} `json:"current"`
			} `json:"breadcrumbs"`
			ZengoIdWorks []struct {
				Id                      string      `json:"id"`
				Title                   string      `json:"title"`
				IllustType              int         `json:"illustType"`
				XRestrict               int         `json:"xRestrict"`
				Restrict                int         `json:"restrict"`
				Sl                      int         `json:"sl"`
				Url                     string      `json:"url"`
				Description             string      `json:"description"`
				Tags                    []string    `json:"tags"`
				UserId                  string      `json:"userId"`
				UserName                string      `json:"userName"`
				Width                   int         `json:"width"`
				Height                  int         `json:"height"`
				PageCount               int         `json:"pageCount"`
				IsBookmarkable          bool        `json:"isBookmarkable"`
				BookmarkData            interface{} `json:"bookmarkData"`
				Alt                     string      `json:"alt"`
				TitleCaptionTranslation struct {
					WorkTitle   interface{} `json:"workTitle"`
					WorkCaption interface{} `json:"workCaption"`
				} `json:"titleCaptionTranslation"`
				CreateDate      time.Time `json:"createDate"`
				UpdateDate      time.Time `json:"updateDate"`
				IsUnlisted      bool      `json:"isUnlisted"`
				IsMasked        bool      `json:"isMasked"`
				AiType          int       `json:"aiType"`
				ProfileImageUrl string    `json:"profileImageUrl"`
			} `json:"zengoIdWorks"`
			ZengoWorkData struct {
				NextWork struct {
					Id    string `json:"id"`
					Title string `json:"title"`
				} `json:"nextWork"`
				PrevWork struct {
					Id    string `json:"id"`
					Title string `json:"title"`
				} `json:"prevWork"`
			} `json:"zengoWorkData"`
		} `json:"noLoginData,omitempty"`
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
