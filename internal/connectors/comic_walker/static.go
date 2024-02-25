package comic_walker

import "time"

type WorkResponse struct {
	Work struct {
		Code       string `json:"code"`
		Id         string `json:"id"`
		Thumbnail  string `json:"thumbnail"`
		BookCover  string `json:"bookCover"`
		Title      string `json:"title"`
		IsOriginal bool   `json:"isOriginal"`
		LabelInfo  struct {
			Name         string `json:"name"`
			IconImageUrl string `json:"iconImageUrl"`
		} `json:"labelInfo"`
		Language string `json:"language"`
		Internal struct {
			DepartmentCode string        `json:"departmentCode"`
			ScrollType     string        `json:"scrollType"`
			LabelNames     []string      `json:"labelNames"`
			FairIds        []interface{} `json:"fairIds"`
		} `json:"internal"`
		Summary string `json:"summary"`
		Genre   struct {
			Code string `json:"code"`
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"genre"`
		SubGenre struct {
			Code string `json:"code"`
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"subGenre"`
		Tags []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"tags"`
		Authors []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
			Role string `json:"role"`
		} `json:"authors"`
		FollowerCount       int    `json:"followerCount"`
		IsNew               bool   `json:"isNew"`
		NextUpdateDateText  string `json:"nextUpdateDateText"`
		IsOneShot           bool   `json:"isOneShot"`
		SerializationStatus string `json:"serializationStatus"`
		RatingLevel         string `json:"ratingLevel"`
	} `json:"work"`
	FirstEpisodes struct {
		Total  int `json:"total"`
		Result []struct {
			Id             string        `json:"id"`
			Code           string        `json:"code"`
			Title          string        `json:"title"`
			SubTitle       string        `json:"subTitle"`
			Thumbnail      string        `json:"thumbnail"`
			DeliveryPeriod time.Time     `json:"deliveryPeriod"`
			IsNew          bool          `json:"isNew"`
			HasRead        bool          `json:"hasRead"`
			Stores         []interface{} `json:"stores"`
			ServiceId      string        `json:"serviceId"`
			Internal       struct {
				EpisodeNo   int    `json:"episodeNo"`
				PageCount   int    `json:"pageCount"`
				Episodetype string `json:"episodetype"`
			} `json:"internal"`
			Type     string `json:"type"`
			IsActive bool   `json:"isActive"`
		} `json:"result"`
	} `json:"firstEpisodes"`
	LatestEpisodes struct {
		Total  int `json:"total"`
		Result []struct {
			Id             string        `json:"id"`
			Code           string        `json:"code"`
			Title          string        `json:"title"`
			SubTitle       string        `json:"subTitle"`
			Thumbnail      string        `json:"thumbnail"`
			DeliveryPeriod time.Time     `json:"deliveryPeriod"`
			IsNew          bool          `json:"isNew"`
			HasRead        bool          `json:"hasRead"`
			Stores         []interface{} `json:"stores"`
			ServiceId      string        `json:"serviceId"`
			Internal       struct {
				EpisodeNo   int    `json:"episodeNo"`
				PageCount   int    `json:"pageCount"`
				Episodetype string `json:"episodetype"`
			} `json:"internal"`
			Type     string `json:"type"`
			IsActive bool   `json:"isActive"`
		} `json:"result"`
	} `json:"latestEpisodes"`
	Comics struct {
		Total  int `json:"total"`
		Result []struct {
			Id        string `json:"id"`
			Title     string `json:"title"`
			Thumbnail string `json:"thumbnail"`
			Release   string `json:"release"`
			Episodes  []struct {
				Id             string        `json:"id"`
				Code           string        `json:"code"`
				Title          string        `json:"title"`
				SubTitle       string        `json:"subTitle"`
				Thumbnail      string        `json:"thumbnail"`
				DeliveryPeriod time.Time     `json:"deliveryPeriod"`
				IsNew          bool          `json:"isNew"`
				HasRead        bool          `json:"hasRead"`
				Stores         []interface{} `json:"stores"`
				ServiceId      string        `json:"serviceId"`
				Internal       struct {
					EpisodeNo   int    `json:"episodeNo"`
					PageCount   int    `json:"pageCount"`
					Episodetype string `json:"episodetype"`
				} `json:"internal"`
				Type     string `json:"type"`
				IsActive bool   `json:"isActive"`
			} `json:"episodes"`
			Stores []struct {
				Code  string `json:"code"`
				Name  string `json:"name"`
				Url   string `json:"url"`
				Image struct {
					Src         string `json:"src"`
					Height      int    `json:"height"`
					Width       int    `json:"width"`
					BlurDataURL string `json:"blurDataURL"`
					BlurWidth   int    `json:"blurWidth"`
					BlurHeight  int    `json:"blurHeight"`
				} `json:"image"`
			} `json:"stores"`
		} `json:"result"`
	} `json:"comics"`
	Promotions   []interface{} `json:"promotions"`
	RelatedBooks struct {
		TotalCount int           `json:"totalCount"`
		Resources  []interface{} `json:"resources"`
	} `json:"relatedBooks"`
	Label struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Color         string `json:"color"`
		IconImageUrl  string `json:"iconImageUrl"`
		Description   string `json:"description"`
		CoverImageUrl string `json:"coverImageUrl"`
		LogoImageUrl  string `json:"logoImageUrl"`
	} `json:"label"`
	Labels []struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Code          string `json:"code"`
		Color         string `json:"color"`
		IconImageUrl  string `json:"iconImageUrl"`
		Description   string `json:"description"`
		CoverImageUrl string `json:"coverImageUrl"`
		LogoImageUrl  string `json:"logoImageUrl"`
	} `json:"labels"`
	LabelWorks []struct {
		Code       string `json:"code"`
		Id         string `json:"id"`
		Thumbnail  string `json:"thumbnail"`
		BookCover  string `json:"bookCover"`
		Title      string `json:"title"`
		IsOriginal bool   `json:"isOriginal"`
		Language   string `json:"language"`
		Internal   struct {
			LabelNames []string `json:"labelNames"`
		} `json:"internal"`
		IsNew   bool `json:"isNew"`
		Episode struct {
			Type  string `json:"type"`
			Code  string `json:"code"`
			Title string `json:"title"`
		} `json:"episode"`
	} `json:"labelWorks"`
	LatestEpisodeId string `json:"latestEpisodeId"`
}

type EpisodeResponse struct {
	Episode struct {
		Id        string `json:"id"`
		Code      string `json:"code"`
		Title     string `json:"title"`
		ServiceId string `json:"serviceId"`
		Thumbnail string `json:"thumbnail"`
		Internal  struct {
			EpisodeNo   int    `json:"episodeNo"`
			PageCount   int    `json:"pageCount"`
			Episodetype string `json:"episodetype"`
			IsLatest    bool   `json:"isLatest"`
		} `json:"internal"`
	} `json:"episode"`
}

type ViewerResponse struct {
	PromotionsEnd []struct {
		Id       string `json:"id"`
		ImageUrl string `json:"imageUrl"`
	} `json:"promotionsEnd"`
	LabelLogo       string    `json:"labelLogo"`
	ScrollDirection string    `json:"scrollDirection"`
	ExpiresAt       time.Time `json:"expiresAt"`
	StartPosition   string    `json:"startPosition"`
	DisplayAds      bool      `json:"displayAds"`
	Manuscripts     []struct {
		DrmMode     string `json:"drmMode"`
		DrmHash     string `json:"drmHash"`
		DrmImageUrl string `json:"drmImageUrl"`
		Page        int    `json:"page"`
	} `json:"manuscripts"`
}

type ExportedCodes struct {
	Work    string
	Episode string
}
