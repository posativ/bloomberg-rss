package domain

type Html struct {
	Props struct {
		PageProps struct {
			Story struct {
				Body Content `json:"body"`
			} `json:"story"`
		} `json:"pageProps"`
	} `json:"props"`
}

type Content struct {
	Type    string    `json:"type"`
	SubType string    `json:"subType"`
	Value   string    `json:"value"`
	Content []Content `json:"content"`
	Data    struct {
		Level  uint   `json:"level"`
		Href   string `json:"href"`
		WebUrl string `json:"webUrl"`
		Photo  struct {
			Src     string `json:"src"`
			Alt     string `json:"alt"`
			Caption string `json:"caption"`
		} `json:"photo"`
		Chart struct {
			Fallback string `json:"fallback"`
			Caption  string `json:"caption"`
		} `json:"chart"`
		Attachment struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Alt         string `json:"alt"`
			Url         string `json:"url"`
		} `json:"attachment"`
	} `json:"data"`
	IFrameData struct {
		Html string `json:"html"`
	} `json:"iframeData"`
	Meta struct {
		Security string `json:"security"`
	} `json:"meta"`
}
