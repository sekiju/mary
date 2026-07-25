package shonenjumpplus

type pageEntry struct {
	Type   string `json:"type"`
	Src    string `json:"src"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type episodeJSON struct {
	ReadableProduct struct {
		Title  string `json:"title"`
		Number int    `json:"number"`
		Series struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"series"`
		PageStructure struct {
			Pages []pageEntry `json:"pages"`
		} `json:"pageStructure"`
	} `json:"readableProduct"`
}

type paginationEntry struct {
	ReadableProductID string `json:"readable_product_id"`
	Title             string `json:"title"`
	ViewerURI         string `json:"viewer_uri"`
}
