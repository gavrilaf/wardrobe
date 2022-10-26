package dto

type FO struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	ContentType string   `json:"content_type"`
	Size        int64    `json:"size"`
	Created     string   `json:"created"`
	Tags        []string `json:"tags"`
}
