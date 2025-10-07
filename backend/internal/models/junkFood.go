package models

// JunkFood represents a junk food entry with a date and type.
type JunkFood struct {
	Date string `json:"date"`
	Type string `json:"type"`
}

type JunkFoodResponse struct {
	Type string `json:"type"`
}
