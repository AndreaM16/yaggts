package parse

import (
	"time"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"strconv"
	"github.com/go-errors/errors"
	"github.com/AndreaM16/yaggts/model"
	"strings"
)

const csvName = "multiTimeline.csv"
const timeLayout = "2006-01-02"
const internalIndex = 5
const weekDaysNumber = 7
var nextDay = [4]int{ 0, 0, 1, -1 }

func GetTrend() ([]model.DbEntry, error) {
	entries := parseCsv(); if len(entries) > 0 {
		return fillMissingDays(entries)
	}
	return []model.DbEntry{}, errors.New("no entries found for last research")
}

// Given a slice of weekly entries, fills missing days between each couple of them
// with a custom average in weekly format only if the latter have different values,
// otherwise they get filled with the previous value.
// - Sunday is the current entry,
// - Wednesday is the average between Sunday and next Sunday
// - Saturday is the average between Wednesday and next Sunday
// - Other days are filled using two indexes, j = 1 and k = 5
//     k = 5; for j = 1; j <= 2; j ++
//       elm[j] = avg(elm[j-1], elm[3])
//       elm[k] = avg(elm[k+1], elm[3]); k--
// Converts each date time.Time in string format yy-mm-dd for Cassandra.
func fillMissingDays(entries []model.CsvEntry) ([]model.DbEntry, error) {
	entriesNumber := len(entries)
	if entriesNumber == 0 {
		return []model.DbEntry{}, errors.New("no entries were appended")
	}
	var trend []model.CsvEntry
	for i := 0; i <= entriesNumber - 2; i++ {
		tmpTrend := make([]model.CsvEntry, 7)
		tmpTrend[0] = entries[i]
		if entries[i].Value != entries[i+1].Value {
			var k = internalIndex
			tmpTrend[3] = model.CsvEntry{
				Value: average(entries[i].Value, entries[i+1].Value),
				Date: entries[i].Date.AddDate(nextDay[0], nextDay[1], nextDay[2] + 2),
			}
			tmpTrend[6] = model.CsvEntry{
				Value: average(entries[i+1].Value, tmpTrend[3].Value),
				Date: entries[i+1].Date.AddDate(nextDay[0], nextDay[1], nextDay[3]),
			}
			for j := 1; j <= 2; j++ {
				tmpTrend[j] = model.CsvEntry{
					Value: average(tmpTrend[j-1].Value, tmpTrend[3].Value),
					Date: tmpTrend[j-1].Date.AddDate(nextDay[0], nextDay[1], nextDay[2]),
				}
				tmpTrend[k] = model.CsvEntry{
					Value: average(tmpTrend[k+1].Value, tmpTrend[3].Value),
					Date:  tmpTrend[k+1].Date.AddDate(nextDay[0], nextDay[1], nextDay[3]),
				}
				k--
			}
		} else {
			tmpTrend[0] = model.CsvEntry{Value: entries[i].Value, Date: entries[i].Date.AddDate(nextDay[0], nextDay[1], nextDay[2])}
			for k := 1; k < weekDaysNumber - 1; k++ {
				tmpTrend[k] = model.CsvEntry{ Value: tmpTrend[k-1].Value, Date: tmpTrend[k-1].Date.AddDate(nextDay[0], nextDay[1], nextDay[2])}
			}
		}
		trend = append(trend, tmpTrend...)
	}; if len(trend) == 0 {
		return []model.DbEntry{}, errors.New("trend is empty")
	}
	trend = append(trend, model.CsvEntry{ Value: entries[len(entries) - 1].Value, Date:  entries[len(entries) - 1].Date })
	finalTrend := []model.DbEntry{}
	for _, t := range trend {
		finalTrend = append(finalTrend, dateToString(t))
	}
	return finalTrend, nil
}

// Returns a date in string format yy-mm-dd
func dateToString(entry model.CsvEntry) model.DbEntry {
	dateEntries := make([]string, 3)
	dateEntries[0] = strconv.Itoa(entry.Date.Year())
	dateEntries[1] = strconv.Itoa(int(entry.Date.Month())); if len(dateEntries[1]) == 1 { dateEntries[1] = "0" + dateEntries[1] }
	dateEntries[2] = strconv.Itoa(entry.Date.Day()); if len(dateEntries[2]) == 1 { dateEntries[2] = "0" + dateEntries[2] }
	m := model.DbEntry{ Value: entry.Value, Date: strings.Join(dateEntries, "-") }
	return m
}

// Returns the average of two float64.
func average(f float64, s float64) float64 {
	return (float64(f) + float64(s)) / float64(2)
}

// Parses a CSV skipping first 3 rows. It returns a slice of csv entries with Date time.Time and Value float64.
// Format is only weekly.
func parseCsv() []model.CsvEntry {
	csvFile, csvFileErr := os.Open(csvName); if csvFileErr != nil {
		return []model.CsvEntry{}
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var trend []model.CsvEntry
	var index = 0
	for {
		line, readerErr := reader.Read(); if readerErr == io.EOF {
			break
		}
		if index > 1 {
			trend = append(trend, model.CsvEntry{
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