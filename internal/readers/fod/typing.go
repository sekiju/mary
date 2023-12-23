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
			Links       interface{} `json:"links,omitempty"`
			Title       string      `json:"title"`
			Author      string      `json:"author"`
			S3Key       string      `json:"s3_key"`
			Direction   string      `json:"direction"`
			HasCover    bool        `json:"has_cover"`
			PageCount   int         `json:"page_count"`
			Navigations []struct {
				Page   int         `json:"page"`
				Text   string      `json:"text"`
				Anchor interface{} `json:"anchor"`
			} `json:"navigations"`
			ImagedReflow bool `json:"imaged_reflow"`
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
			Links       interface{} `json:"links,omitempty"`
			Title       string      `json:"title"`
			Author      string      `json:"author"`
			S3Key       string      `json:"s3_key"`
			Direction   string      `json:"direction"`
			HasCover    bool        `json:"has_cover"`
			PageCount   int         `json:"page_count"`
			Navigations []struct {
				Page   int         `json:"page"`
				Text   string      `json:"text"`
				Anchor interface{} `json:"anchor"`
			} `json:"navigations"`
			ImagedReflow    bool   `json:"imaged_reflow"`
			PageArrangement string `json:"page_arrangement"`
		} `json:"data_for_browser"`
		KeysForBrowser []string `json:"keys_for_browser"`
	} `json:"guardian_info_all"`
}
