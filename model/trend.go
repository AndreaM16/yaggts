package model

type TrendEntry struct {
	Value float64 `json:"value"`
	Date string `json:"date"`
}


type Trend struct {
	Query string `json:"query,omitempty"`
	Trend []TrendEntry
}