package model

type DbEntry struct {
	Value float64 `json:"value"`
	Date string `json:"date"`
}


type Trend struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Trend []DbEntry
}