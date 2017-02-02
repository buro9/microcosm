package models

// Link is a hyperlink
type Link struct {
	Rel   string `json:"rel,omitempty"` // REST
	Href  string `json:"href"`
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"` // HTML
}
