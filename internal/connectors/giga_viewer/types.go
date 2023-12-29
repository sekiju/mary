package giga_viewer

import "time"

type EpisodeResponse struct {
	ReadableProduct struct {
		FinishReadingNotificationUri interface{} `json:"finishReadingNotificationUri"`
		HasPurchased                 bool        `json:"hasPurchased"`
		Id                           string      `json:"id"`
		ImageUrisDigest              string      `json:"imageUrisDigest"`
		IsPublic                     bool        `json:"isPublic"`
		NextReadableProductUri       interface{} `json:"nextReadableProductUri"`
		Number                       int         `json:"number"`
		PageStructure                struct {
			ChoJuGiga string `json:"choJuGiga"`
			Pages     []struct {
				Buttons []struct {
					Type string `json:"type"`
					Uri  string `json:"uri"`
				} `json:"buttons,omitempty"`
				ContentStart string `json:"contentStart,omitempty"`
				Height       int    `json:"height,omitempty"`
				Src          string `json:"src,omitempty"`
				Type         string `json:"type"`
				Width        int    `json:"width,omitempty"`
				ContentEnd   string `json:"contentEnd,omitempty"`
				LinkPosition string `json:"linkPosition,omitempty"`
			} `json:"pages"`
			ReadingDirection string `json:"readingDirection"`
			StartPosition    string `json:"startPosition"`
		} `json:"pageStructure"`
		Permalink                               string      `json:"permalink"`
		PointGettableEpisodeWhenCompleteReading interface{} `json:"pointGettableEpisodeWhenCompleteReading"`
		PrevReadableProductUri                  string      `json:"prevReadableProductUri"`
		PublishedAt                             time.Time   `json:"publishedAt"`
		Series                                  struct {
			Id           string `json:"id"`
			ThumbnailUri string `json:"thumbnailUri"`
			Title        string `json:"title"`
		} `json:"series"`
		Title    string      `json:"title"`
		Toc      interface{} `json:"toc"`
		TypeName string      `json:"typeName"`
	} `json:"readableProduct"`
}
