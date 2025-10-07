package models

// Activity represents an activity item.
type Activity struct {
	Date string `json:"date"`
	Type string `json:"type"`
	// ElapsedTime in seconds (0 if not set)
	ElapsedTime int `json:"elapsedTime"`
}

type ActivityResponse struct {
	Type        string `json:"type"`
	ElapsedTime int    `json:"elapsedTime"`
}
