package giga_viewer

import "time"

type EpisodeResponse struct {
	ReadableProduct struct {
		HasPurchased  bool   `json:"hasPurchased"`
		Id            string `json:"id"`
		IsPublic      bool   `json:"isPublic"`
		Number        int    `json:"number"`
		PageStructure *struct {
			ChoJuGiga        string        `json:"choJuGiga"`
			Pages            []EpisodePage `json:"pages"`
			ReadingDirection string        `json:"readingDirection"`
			StartPosition    string        `json:"startPosition"`
		} `json:"pageStructure,omitempty"`
		PublishedAt time.Time `json:"publishedAt"`
		Series      struct {
			Id           string `json:"id"`
			ThumbnailUri string `json:"thumbnailUri"`
			Title        string `json:"title"`
		} `json:"series"`
		Title    string `json:"title"`
		TypeName string `json:"typeName"`
	} `json:"readableProduct"`
}

type EpisodePage struct {
	Src  string `json:"src,omitempty"`
	Type string `json:"type"`
}
