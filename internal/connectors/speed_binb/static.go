package speed_binb

import "image"

type DescrambleView struct {
	Width     int
	Height    int
	Transfers []DescrambleTransfer
}

type DescrambleTransfer struct {
	Index  int
	Coords []DescrambleCord
}

type Piece struct {
	x, y, w, h int
}

type Yt struct {
	ndx, ndy int
	piece    *[]Piece
}

type SpeedBinbDecoder interface {
	vt() bool
	bt(t image.Rectangle) bool
	dt(t image.Rectangle) image.Rectangle
	gt(t image.Rectangle) []DescrambleCord
}

type DescrambleCord struct {
	XSrc   int
	YSrc   int
	Width  int
	Height int
	XDest  int
	YDest  int
}

// json...

type Ptimg struct {
	PtimgVersion int `json:"ptimg-version"`
	Resources    struct {
		I struct {
			Src    string `json:"src"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"i"`
	} `json:"resources"`
	Views []PtimgView `json:"views"`
}

type PtimgView struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Coords []string `json:"coords"`
}

type BibGetCntntInfo struct {
	Result     int    `json:"result"`
	ShopUserID string `json:"ShopUserID"`
	Eurl       string `json:"eurl"`
	Items      []struct {
		ContentID      string `json:"ContentID"`
		ContentsServer string `json:"ContentsServer"`
		ServerType     int    `json:"ServerType"`
		Authors        []struct {
			Name string `json:"Name"`
			Ruby string `json:"Ruby"`
		} `json:"Authors"`
		Publisher                      string   `json:"Publisher"`
		PublisherRuby                  string   `json:"PublisherRuby"`
		Title                          string   `json:"Title"`
		TitleRuby                      string   `json:"TitleRuby"`
		SubTitle                       string   `json:"SubTitle"`
		Categories                     []string `json:"Categories"`
		Abstract                       string   `json:"Abstract"`
		Description                    string   `json:"Description"`
		ContentType                    string   `json:"ContentType"`
		TermForRead                    string   `json:"TermForRead"`
		LastPageDim                    int      `json:"LastPageDim"`
		LastPageColor                  string   `json:"LastPageColor"`
		LastPageURL                    string   `json:"LastPageURL"`
		ShopURL                        string   `json:"ShopURL"`
		ThumbnailImageURL              string   `json:"ThumbnailImageURL"`
		IPAddress                      string   `json:"IPAddress"`
		ViewMode                       int      `json:"ViewMode"`
		RegistDate                     string   `json:"RegistDate"`
		Stbl                           string   `json:"stbl"`
		Ttbl                           string   `json:"ttbl"`
		Ptbl                           string   `json:"ptbl"`
		Ctbl                           string   `json:"ctbl"`
		SliderCaption                  int      `json:"SliderCaption"`
		P                              string   `json:"p"`
		ConverterVersion               string   `json:"ConverterVersion"`
		ContentDate                    string   `json:"ContentDate"`
		InlineRecommendPageURL         string   `json:"InlineRecommendPageURL"`
		CmoaPopupRecommendPageURL      string   `json:"CmoaPopupRecommendPageURL"`
		RequestShowOperationTips       int64    `json:"RequestShowOperationTips"`
		OutOfRangeAutoBookmarkBehavior string   `json:"OutOfRangeAutoBookmarkBehavior"`
	} `json:"items"`
}

type BibGetCntntInfo16130 struct {
	Result          int         `json:"result"`
	ShopUserID      interface{} `json:"ShopUserID"`
	Eurl            string      `json:"eurl"`
	Rurl            string      `json:"rurl"`
	BannerCloseTime int         `json:"BannerCloseTime"`
	TerminalType    int         `json:"TerminalType"`
	MenuUpdateTime  string      `json:"MenuUpdateTime"`
	Items           []struct {
		ContentID      string `json:"ContentID"`
		ContentsServer string `json:"ContentsServer"`
		ServerType     int    `json:"ServerType"`
		Authors        []struct {
			Name string      `json:"Name"`
			Ruby interface{} `json:"Ruby"`
			Role interface{} `json:"Role"`
			Path string      `json:"Path"`
		} `json:"Authors"`
		Publisher                      string        `json:"Publisher"`
		PublisherRuby                  string        `json:"PublisherRuby"`
		Title                          string        `json:"Title"`
		ParentTitle                    string        `json:"ParentTitle"`
		EpisodeThumbnail               string        `json:"EpisodeThumbnail"`
		ParentDescription              string        `json:"ParentDescription"`
		TitleRuby                      string        `json:"TitleRuby"`
		SubTitle                       string        `json:"SubTitle"`
		Categories                     []interface{} `json:"Categories"`
		Abstract                       string        `json:"Abstract"`
		Description                    string        `json:"Description"`
		ContentType                    int           `json:"ContentType"`
		MostRecommendUser              interface{}   `json:"MostRecommendUser"`
		LastPageDim                    int           `json:"LastPageDim"`
		LastPageColor                  string        `json:"LastPageColor"`
		LastPageURL                    string        `json:"LastPageURL"`
		LastPageCssURL                 string        `json:"LastPageCssURL"`
		RecommendPageURL               string        `json:"RecommendPageURL"`
		WarningPageURL                 string        `json:"WarningPageURL"`
		WarningPageCssURL              string        `json:"WarningPageCssURL"`
		ShopURL                        string        `json:"ShopURL"`
		OutOfRangeAutoBookmarkBehavior string        `json:"OutOfRangeAutoBookmarkBehavior"`
		ThumbnailImageURL              string        `json:"ThumbnailImageURL"`
		IPAddress                      string        `json:"IPAddress"`
		ViewMode                       int           `json:"ViewMode"`
		Stbl                           string        `json:"stbl"`
		Ttbl                           string        `json:"ttbl"`
		Ptbl                           string        `json:"ptbl"`
		Ctbl                           string        `json:"ctbl"`
		P                              interface{}   `json:"p"`
		SliderCaption                  int           `json:"SliderCaption"`
		FirstBanner                    []struct {
			Index int    `json:"index"`
			Name  string `json:"name"`
			Url   string `json:"url"`
			Img   string `json:"img"`
		} `json:"FirstBanner"`
		Banners struct {
			AfterViewer1 []struct {
				Index int    `json:"index"`
				Name  string `json:"name"`
				Url   string `json:"url"`
				Img   string `json:"img"`
			} `json:"after_viewer_1"`
			AfterViewer2 []struct {
				Index int    `json:"index"`
				Name  string `json:"name"`
				Url   string `json:"url"`
				Img   string `json:"img"`
			} `json:"after_viewer_2"`
			AfterViewer3 []struct {
				Index int    `json:"index"`
				Name  string `json:"name"`
				Url   string `json:"url"`
				Img   string `json:"img"`
			} `json:"after_viewer_3"`
		} `json:"Banners"`
		ShareURL                 string `json:"ShareURL"`
		TwitterHashTag           string `json:"TwitterHashTag"`
		Price                    string `json:"Price"`
		UseGuidePage             int    `json:"UseGuidePage"`
		RequestShowOperationTips string `json:"RequestShowOperationTips"`
		RelativeUrlRoot          string `json:"RelativeUrlRoot"`
		ConverterVersion         string `json:"ConverterVersion"`
		InlineRecommendPageURL   []struct {
			Url string `json:"Url"`
		} `json:"InlineRecommendPageURL"`
		ParentPath      string      `json:"ParentPath"`
		Customer        interface{} `json:"Customer"`
		PrevEpisode     interface{} `json:"PrevEpisode"`
		NextEpisode     string      `json:"NextEpisode"`
		FavoritePageUrl string      `json:"FavoritePageUrl"`
	} `json:"items"`
}

type SbcGetCntnt struct {
	SBCVersion       string  `json:"SBCVersion"`
	Result           int     `json:"result"`
	Ttx              string  `json:"ttx"`
	Prop             string  `json:"prop"`
	AddressList      [][]int `json:"AddressList"`
	ConverterType    string  `json:"ConverterType"`
	ConverterVersion string  `json:"ConverterVersion"`
	ContentDate      string  `json:"ContentDate"`
	NecImageSize     int     `json:"NecImageSize"`
	NecImageCnt      int     `json:"NecImageCnt"`
	SmlImageSizeHQ   int     `json:"SmlImageSizeHQ"`
	SmlImageSizeLQ   int     `json:"SmlImageSizeLQ"`
	SmlImageCnt      int     `json:"SmlImageCnt"`
	IsTateyomi       bool    `json:"IsTateyomi"`
	OffDataType      int     `json:"OffDataType"`
}
