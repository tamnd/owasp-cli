package owasp

// Sheet is a single OWASP Cheat Sheet.
type Sheet struct {
	Rank int    `json:"rank"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Raw  string `json:"raw"`
	Size int    `json:"size"`
}

// wireFile is the response shape from the GitHub Contents API.
type wireFile struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
	HTMLURL     string `json:"html_url"`
	Size        int    `json:"size"`
	Type        string `json:"type"`
}
