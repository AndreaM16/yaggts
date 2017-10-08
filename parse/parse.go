package parse

import (
	"time"
	"os"
	"fmt"
	"encoding/csv"
	"bufio"
	"io"
	"strconv"
	"github.com/go-errors/errors"
)

const csvName = "multiTimeline.csv"
const timeLayout = "2006-01-02"
const internalIndex = 5
const weekDaysNumber = 7
var nextDay = [4]int{ 0, 0, 1, -1 }

type CsvEntry struct {
	Date time.Time `json:"date"`
	Value float64 `json:"value"`
}

func GetTrend() ([]CsvEntry, error) {
	entries := parseCsv()
	entriesNumber := len(entries)
	if entriesNumber == 0 {
		return entries, errors.New("No entries were appended.")
	}
	var trend []CsvEntry

	for i := 0; i <= entriesNumber - 2; i++ {
		tmpTrend := make([]CsvEntry, 7)
		tmpTrend[0] = entries[i]
		if entries[i].Value != entries[i+1].Value {
			var decrementIndex = internalIndex
			tmpTrend[3] = CsvEntry{
				Value: average(entries[i].Value, entries[i+1].Value),
				Date: entries[i].Date.AddDate(nextDay[0], nextDay[1], nextDay[2] + 2),
			}
			tmpTrend[6] = CsvEntry{
				Value: average(entries[i+1].Value, tmpTrend[3].Value),
				Date: entries[i+1].Date.AddDate(nextDay[0], nextDay[1], nextDay[3]),
			}
			for j := 1; j <= 2; j++ {
				tmpTrend[j] = CsvEntry{
					Value: average(tmpTrend[j-1].Value, tmpTrend[3].Value),
					Date: tmpTrend[j-1].Date.AddDate(nextDay[0], nextDay[1], nextDay[2]),
				}
				tmpTrend[decrementIndex] = CsvEntry{
					Value: average(tmpTrend[decrementIndex + 1].Value, tmpTrend[3].Value),
					Date: tmpTrend[decrementIndex + 1].Date.AddDate(nextDay[0], nextDay[1], nextDay[3]),
				}
				decrementIndex--
			}
		} else {
			tmpTrend[0] = CsvEntry{Value: entries[i].Value, Date: entries[i].Date.AddDate(nextDay[0], nextDay[1], nextDay[2])}
			for k := 1; k <= weekDaysNumber - 1; k++ {
				tmpTrend[k] = CsvEntry{ Value: tmpTrend[k-1].Value, Date: tmpTrend[k-1].Date.AddDate(nextDay[0], nextDay[1], nextDay[2])}
			}
		}
		trend = append(trend, tmpTrend...)
	}; if len(trend) == 0 {
		return trend, errors.New("Trend is empty")
	}
	trend = append(trend, CsvEntry{ Value: entries[len(entries) - 1].Value, Date:  entries[len(entries) - 1].Date })
	for _, k := range trend {
		fmt.Println(k.Value, k.Date)
	}
	return trend, nil
}

func average(f float64, s float64) float64 {
	return (float64(f) + float64(s)) / float64(2)
}

func parseCsv() []CsvEntry {
	csvFile, csvFileErr := os.Open(csvName); if csvFileErr != nil {
		fmt.Println("Unable to open csv.")
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var trend []CsvEntry
	var index = 0
	for {
		line, readerErr := reader.Read(); if readerErr == io.EOF {
			break
		}
		if index > 1 {
			trend = append(trend, CsvEntry{
				Date: func (val string) time.Time {
					t, _ := time.Parse(timeLayout, val); return t
				}(line[0]),
				Value: func (val string) float64 {
					v, e := strconv.ParseFloat(val, 64); if e == nil {
						return v
					}
					return 0
				}(line[1]),
			})
		}
		index++
	}; return trend
}