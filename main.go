package main

import (
	"fmt"
	"strings"
	"os"
	"github.com/andream16/yaggts/model"
	executor "github.com/andream16/yaggts/exec"
	parser "github.com/andream16/yaggts/parse"
)

var query string

func init() {
	// Read program arguments and join them in a string
	query = strings.Join(os.Args[1:], " ")
}

func main() {
	fmt.Println("Starting yaggts . . .")
	results := func() bool {
		return executor.DownloadCSV(query)
	}()
	if !results {
		fmt.Println("Unable to download csv.")
		return
	}
	fmt.Println("Got the csv.")
	trend, trendErr := parser.GetTrend()
	var trendEntry model.Trend
	if trendErr == nil && len(trend) > 0 {
		fmt.Println("Got a valid trend.")
		for _, k := range trend {
			fmt.Println(fmt.Sprintf("Date %s ; Value %v", k.Date, k.Value))
			tuple := model.TrendEntry{
				Date: k.Date,
				Value: k.Value,
			}
			trendEntry.Trend = append(trendEntry.Trend, tuple)
		}
	} else {
		fmt.Println(trendErr.Error())
	}
}