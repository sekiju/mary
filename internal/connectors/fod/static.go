package fod

import "time"

type LicenceKeyResponse struct {
	ServerStatus struct {
		AccessTime time.Time `json:"access_time"`
		AppVersion string    `json:"app_version"`
		ResultCode int       `json:"result_code"`
	} `json:"server_status"`
	GuardianInfoForBrowser struct {
		GUARDIANSERVER        string `json:"GUARDIAN_SERVER"`
		ADDITIONALQUERYSTRING string `json:"ADDITIONAL_QUERY_STRING"`
		BookmarkIdPrefix      string `json:"bookmarkIdPrefix"`
		BookData              struct {
			Extra []interface{} `json:"extra"`
			Links interface{}   `json:"links"`
			Spine []struct {
				Path   string `json:"path"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Linear bool   `json:"linear"`
			} `json:"spine"`
			Title       string `json:"title"`
			Author      string `json:"author"`
			S3Key       string `json:"s3_key"`
			Version     int    `json:"version"`
			Direction   string `json:"direction"`
			HasCover    bool   `json:"has_cover"`
			PageCount   int    `json:"page_count"`
			Navigations []struct {
				Text       string      `json:"text"`
				Anchor     interface{} `json:"anchor"`
				NestLevel  int         `json:"nest_level"`
				SpineIndex int         `json:"spine_index"`
			} `json:"navigations"`
			ContentType  string `json:"content_type"`
			ImagedReflow bool   `json:"imaged_reflow"`
		} `json:"book_data"`
		PagesData struct {
			Start           int      `json:"start"`
			End             int      `json:"end"`
			Bookmark        int      `json:"bookmark"`
			Tableofcontents int      `json:"tableofcontents"`
			Keys            []string `json:"keys"`
		} `json:"pages_data"`
		HasFrontCover                bool `json:"has_front_cover"`
		HideLogo                     bool `json:"hide_logo"`
		CloseBtn                     bool `json:"close_btn"`
		BackBtn                      bool `json:"back_btn"`
		ClearAutobookmarkAtBackcover bool `json:"clear_autobookmark_at_backcover"`
		UseBtnTitle                  bool `json:"use_btn_title"`
		BottomToolbar                bool `json:"bottom_toolbar"`
		FitWidthOnMobileLandscape    bool `json:"fitWidthOnMobileLandscape"`
		OnlyPortraitImageOnReflow    bool `json:"onlyPortraitImageOnReflow"`
		WebtoonFlg                   bool `json:"webtoon_flg"`
	} `json:"guardian_info_for_browser"`
	GuardianInfoAll struct {
		ContentId      string `json:"content_id"`
		NativeS3Key    string `json:"native_s3_key"`
		ImagedReflow   bool   `json:"imaged_reflow"`
		DataForBrowser struct {
			Extra []interface{} `json:"extra"`
			Links interface{}   `json:"links"`
			Spine []struct {
				Path   string `json:"path"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Linear bool   `json:"linear"`
			} `json:"spine"`
			Title       string `json:"title"`
			Author      string `json:"author"`
			S3Key       string `json:"s3_key"`
			Version     int    `json:"version"`
			Direction   string `json:"direction"`
			HasCover    bool   `json:"has_cover"`
			PageCount   int    `json:"page_count"`
			Navigations []struct {
				Text       string      `json:"text"`
				Anchor     interface{} `json:"anchor"`
				NestLevel  int         `json:"nest_level"`
				SpineIndex int         `json:"spine_index"`
			} `json:"navigations"`
			ContentType     string `json:"content_type"`
			ImagedReflow    bool   `json:"imaged_reflow"`
			PageArrangement string `json:"page_arrangement"`
		} `json:"data_for_browser"`
		KeysForBrowser []string `json:"keys_for_browser"`
	} `json:"guardian_info_all"`
}
type DetailResponse struct {
	ServerStatus struct {
		AccessTime time.Time `json:"access_time"`
		AppVersion string    `json:"app_version"`
		ResultCode int       `json:"result_code"`
	} `json:"server_status"`
	User struct {
		IsPremiumUser bool `json:"is_premium_user"`
	} `json:"user"`
	BookDetail struct {
		BookId              string      `json:"book_id"`
		BookStatus          int         `json:"book_status"`
		BookStartDate       string      `json:"book_start_date"`
		BookName            string      `json:"book_name"`
		BookReviewLong      string      `json:"book_review_long"`
		EpisodeId           string      `json:"episode_id"`
		EpisodeCount        int         `json:"episode_count"`
		EpisodeName         string      `json:"episode_name"`
		EpisodeStatus       int         `json:"episode_status"`
		EpisodeLicenseEnd   time.Time   `json:"episode_license_end"`
		EpisodeReleaseStart time.Time   `json:"episode_release_start"`
		EpisodeCloseEnd     time.Time   `json:"episode_close_end"`
		IsSample            bool        `json:"is_sample"`
		IsPurchased         bool        `json:"is_purchased"`
		EpisodeReviewLong   string      `json:"episode_review_long"`
		CanReadFree         bool        `json:"can_read_free"`
		Likes               int         `json:"likes"`
		IsLiked             bool        `json:"is_liked"`
		WebtoonFlag         bool        `json:"webtoon_flag"`
		ReadAt              interface{} `json:"read_at"`
		IsPremiumFree       bool        `json:"is_premium_free"`
		RegularPrice        int         `json:"regular_price"`
		EpisodePriceStart   string      `json:"episode_price_start"`
		EpisodePriceEnd     string      `json:"episode_price_end"`
		CashbackRate        float32     `json:"cashback_rate"`
		CashbackCloseDate   string      `json:"cashback_close_date"`
		CashbackPoint       int         `json:"cashback_point"`
		DiscountRate        int         `json:"discount_rate"`
		DiscountEndDate     string      `json:"discount_end_date"`
		DiscountedPrice     int         `json:"discounted_price"`
		IsWishlist          bool        `json:"is_wishlist"`
		DedicatedAnnotation string      `json:"dedicated_annotation"`
		VisibleReadButton   bool        `json:"visible_read_button"`
		Thumbnail           string      `json:"thumbnail"`
		IsFree              bool        `json:"is_free"`
		IsNew               bool        `json:"is_new"`
		IsConclusion        bool        `json:"is_conclusion"`
		IsLimited           bool        `json:"is_limited"`
		IsNewIssue          bool        `json:"is_new_issue"`
		IsSequel            bool        `json:"is_sequel"`
		Authors             []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"authors"`
		Publishers []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"publishers"`
		Genres []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
		Magazines []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"magazines"`
		SubGenres []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"sub_genres"`
		IsColor      bool `json:"is_color"`
		IsPhotoAlbum bool `json:"is_photo_album"`
		IsNovel      bool `json:"is_novel"`
		UserBook     struct {
			UserShelfId      interface{} `json:"user_shelf_id"`
			UserEpisodeCount int         `json:"user_episode_count"`
			BookNameKana     string      `json:"book_name_kana"`
			SeriesCount      int         `json:"series_count"`
			ReadAt           interface{} `json:"read_at"`
			SummaryDate      interface{} `json:"summary_date"`
			AuthorNameKana   string      `json:"author_name_kana"`
		} `json:"user_book"`
	} `json:"book_detail"`
	BookSeries []struct {
		BookId            string    `json:"book_id"`
		BookName          string    `json:"book_name"`
		EpisodeId         string    `json:"episode_id"`
		EpisodeCount      int       `json:"episode_count"`
		EpisodeLicenseEnd time.Time `json:"episode_license_end"`
		IsSample          bool      `json:"is_sample,omitempty"`
		IsPurchased       bool      `json:"is_purchased"`
		Likes             int       `json:"likes"`
		IsLiked           bool      `json:"is_liked"`
		IsPremiumFree     bool      `json:"is_premium_free"`
		RegularPrice      int       `json:"regular_price"`
		EpisodePriceStart string    `json:"episode_price_start"`
		EpisodePriceEnd   string    `json:"episode_price_end"`
		CashbackRate      float32   `json:"cashback_rate"`
		CashbackCloseDate string    `json:"cashback_close_date"`
		CashbackPoint     int       `json:"cashback_point"`
		DiscountRate      int       `json:"discount_rate"`
		DiscountEndDate   string    `json:"discount_end_date"`
		DiscountedPrice   int       `json:"discounted_price"`
		IsWishlist        bool      `json:"is_wishlist"`
		Thumbnail         string    `json:"thumbnail"`
		IsFree            bool      `json:"is_free"`
		IsNew             bool      `json:"is_new"`
		IsConclusion      bool      `json:"is_conclusion"`
		IsLimited         bool      `json:"is_limited"`
		IsNewIssue        bool      `json:"is_new_issue"`
		IsSequel          bool      `json:"is_sequel"`
		Authors           []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"authors"`
		Publishers []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"publishers"`
		Genres []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
	} `json:"book_series"`
}

type ServerStatusResponse struct {
	ServerStatus struct {
		AccessTime time.Time `json:"access_time"`
		AppVersion string    `json:"app_version"`
		ResultCode int       `json:"result_code"`
	} `json:"server_status"`
}

type BookCredentialsRequest struct {
	BookID    string `json:"book_id"`
	EpisodeID string `json:"episode_id"`
}
