package upload

type UploadImageResponse struct {
	OK   bool     `json:"ok"`
	URL  string   `json:"url,omitempty"`
	URLs []string `json:"urls,omitempty"`
}
