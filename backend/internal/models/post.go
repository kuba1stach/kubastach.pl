package models

// Media represents a media attachment within a post.
type Media struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Post represents a blog / content post item.
type Post struct {
	Content string  `json:"content"`
	Date    string  `json:"date"`
	Media   []Media `json:"media"`
}
