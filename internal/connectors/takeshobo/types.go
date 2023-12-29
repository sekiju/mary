package takeshobo

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
