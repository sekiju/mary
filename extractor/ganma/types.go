package ganma

type createAccountResponse struct {
	Success bool `json:"success"`
	Root    struct {
		Id       string `json:"id"`
		Password string `json:"password"`
	} `json:"root"`
}

type magazineResult struct {
	Success bool `json:"success"`
	Root    struct {
		Id                        string `json:"id"`
		Alias                     string `json:"alias"`
		Title                     string `json:"title"`
		Description               string `json:"description"`
		Overview                  string `json:"overview"`
		RectangleWithLogoImageURL string `json:"rectangleWithLogoImageURL"`
		DistributionLabel         string `json:"distributionLabel"`
		StoryReleaseStatus        string `json:"storyReleaseStatus"`
		PublicLatestStoryNumber   int    `json:"publicLatestStoryNumber"`
		Upcoming                  struct {
			Title         string `json:"title"`
			ScheduledDate int64  `json:"scheduledDate"`
		} `json:"upcoming"`
		Author struct {
			Name            string `json:"name"`
			ProfileText     string `json:"profileText"`
			ProfileImageURL string `json:"profileImageURL"`
			Link            struct {
				Text string `json:"text"`
				Url  string `json:"url"`
			} `json:"link"`
		} `json:"author"`
		IsSeriesBind bool `json:"isSeriesBind"`
		Items        []struct {
			StoryId           string `json:"storyId"`
			Title             string `json:"title"`
			SeriesTitle       string `json:"seriesTitle"`
			ThumbnailImageURL string `json:"thumbnailImageURL"`
			Number            int    `json:"number,omitempty"`
			Kind              string `json:"kind"`
			ReleaseStart      int64  `json:"releaseStart"`
			HeartCount        int    `json:"heartCount"`
			DisableCM         bool   `json:"disableCM"`
			HasExchange       bool   `json:"hasExchange"`
			Subtitle          string `json:"subtitle,omitempty"`
		} `json:"items"`
		Recommendations []interface{} `json:"recommendations"`
		RelatedLink     struct {
			Label string `json:"label"`
			Items []struct {
				Name       string `json:"name"`
				ImageURL   string `json:"imageURL"`
				Transition struct {
					Way            string `json:"way"`
					DestinationURL string `json:"destinationURL"`
				} `json:"transition"`
			} `json:"items"`
		} `json:"relatedLink"`
		CanSupport              bool          `json:"canSupport"`
		CanAcceptFanLetter      bool          `json:"canAcceptFanLetter"`
		FirstViewAdvertisements []interface{} `json:"firstViewAdvertisements"`
		FooterAdvertisements    []interface{} `json:"footerAdvertisements"`
		Tags                    []struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"tags"`
		HeartCount         int      `json:"heartCount"`
		BookmarkCount      int      `json:"bookmarkCount"`
		HighlightImageURLs []string `json:"highlightImageURLs"`
		IsGTOON            bool     `json:"isGTOON"`
	} `json:"root"`
}

type readerResult struct {
	Errors []struct {
		Message    string `json:"message"`
		Extensions struct {
			Code string `json:"code"`
		} `json:"extensions"`
	} `json:"errors"`
	Data struct {
		Magazine struct {
			MagazineId string `json:"magazineId"`
			Title      string `json:"title"`
			Alias      string `json:"alias"`
			Overview   string `json:"overview"`
			Author     struct {
				PenName         string `json:"penName"`
				ProfileImageURL string `json:"profileImageURL"`
			} `json:"author"`
			AuthorName                string `json:"authorName"`
			RectangleWithLogoImageURL string `json:"rectangleWithLogoImageURL"`
			SquareImageURL            string `json:"squareImageURL"`
			IsSellByStoryMagazine     bool   `json:"isSellByStoryMagazine"`
			ShareText                 string `json:"shareText"`
			CanSupport                bool   `json:"canSupport"`
			Upcoming                  struct {
				Title     string      `json:"title"`
				Subtitle  interface{} `json:"subtitle"`
				ReleaseAt int64       `json:"releaseAt"`
			} `json:"upcoming"`
			StoryContents struct {
				Typename  string `json:"__typename"`
				StoryInfo struct {
					StoryId                 string      `json:"storyId"`
					Title                   string      `json:"title"`
					Subtitle                interface{} `json:"subtitle"`
					AppealMessage           interface{} `json:"appealMessage"`
					IsVerticalOnly          bool        `json:"isVerticalOnly"`
					ContentsAccessCondition struct {
						Typename  string `json:"__typename"`
						DisableCM bool   `json:"disableCM"`
					} `json:"contentsAccessCondition"`
					NextStoryInfo struct {
						Typename                string      `json:"__typename"`
						StoryId                 string      `json:"storyId"`
						Title                   string      `json:"title"`
						Subtitle                interface{} `json:"subtitle"`
						AppealMessage           interface{} `json:"appealMessage"`
						IsPurchased             bool        `json:"isPurchased"`
						StoryThumbnailURLs      []string    `json:"storyThumbnailURLs"`
						ContentsAccessCondition struct {
							Typename  string `json:"__typename"`
							DisableCM bool   `json:"disableCM"`
						} `json:"contentsAccessCondition"`
					} `json:"nextStoryInfo"`
					PreviousStoryInfo struct {
						Typename                string      `json:"__typename"`
						StoryId                 string      `json:"storyId"`
						Title                   string      `json:"title"`
						Subtitle                interface{} `json:"subtitle"`
						AppealMessage           interface{} `json:"appealMessage"`
						IsPurchased             bool        `json:"isPurchased"`
						StoryThumbnailURLs      []string    `json:"storyThumbnailURLs"`
						ContentsAccessCondition struct {
							Typename  string `json:"__typename"`
							DisableCM bool   `json:"disableCM"`
						} `json:"contentsAccessCondition"`
					} `json:"previousStoryInfo"`
					GeneralUpcoming interface{} `json:"generalUpcoming"`
					HeartCount      int         `json:"heartCount"`
					ReleaseForFree  interface{} `json:"releaseForFree"`
				} `json:"storyInfo"`
				PageImages struct {
					SecretKey        interface{} `json:"secretKey"`
					PageImageBaseURL string      `json:"pageImageBaseURL"`
					PageImageSign    string      `json:"pageImageSign"`
					PageCount        int         `json:"pageCount"`
				} `json:"pageImages"`
				StoryEndImage           interface{} `json:"storyEndImage"`
				StoryEndImageOnMagazine interface{} `json:"storyEndImageOnMagazine"`
				Afterword               struct {
					Text     interface{} `json:"text"`
					ImageURL string      `json:"imageURL"`
				} `json:"afterword"`
				StoryEndAd1 interface{} `json:"storyEndAd1"`
				StoryEndAd2 interface{} `json:"storyEndAd2"`
				OverlayAd   interface{} `json:"overlayAd"`
				Exchange    struct {
					CoverImageURL string `json:"coverImageURL"`
				} `json:"exchange"`
			} `json:"storyContents"`
		} `json:"magazine"`
	} `json:"data"`
}
