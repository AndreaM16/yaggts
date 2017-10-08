package model

import "time"

type CsvEntry struct {
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
}