package main

import (
	"fmt"
	"log"
	executor "github.com/AndreaM16/yaggts/exec"
	parser "github.com/AndreaM16/yaggts/parse"
	//"github.com/andream16/yaggts/util"
	"github.com/AndreaM16/yaggts/util"
	"encoding/json"
	"net/http"
	"bytes"
	"errors"
	"github.com/AndreaM16/yaggts/model"
)

func main() {
	fmt.Println("Starting yaggts . . .")

	manufacturers, manufacturersError := util.GetManufacturers(); if manufacturersError != nil {
		panic(manufacturersError)
	}
	for _, manufacturer := range manufacturers.Manufacturers {
		results := func() bool {
			return executor.DownloadCSV(manufacturer.Manufacturer)
		}()
		if !results {
			log.Fatal("Unable to download csv.")
		}
		fmt.Println("Got the csv.")
		trend, trendErr := parser.GetTrend()
		var trendEntry model.Trend
		if trendErr == nil && len(trend) > 0 {
			fmt.Println("Got a valid trend.")
			for _, k := range trend {
				fmt.Println(k.Value, k.Date)
				var dbEntry model.DbEntry
				trendEntry.Manufacturer = manufacturer.Manufacturer
				dbEntry.Date = k.Date
				dbEntry.Value = k.Value
				trendEntry.Trend = append(trendEntry.Trend, dbEntry)
			}
			postTrend(trendEntry)
			//util.DeleteCsv()
		} else {
			fmt.Println(trendErr.Error())
		}
	}
}

func postTrend(trendEntry model.Trend) error {
	fmt.Println(fmt.Sprintf("Posting new trend entry for manufacturer %s . . .", trendEntry.Manufacturer))
	body, bodyError := json.Marshal(trendEntry); if bodyError != nil {
		fmt.Println(fmt.Sprintf("Unable to marshal new trend entry for manufacturer %s, got error: %s", trendEntry.Manufacturer, bodyError.Error()))
		return bodyError
	}
	response, requestErr := http.Post("http://localhost:8080/trend", "application/json", bytes.NewBuffer(body)); if requestErr != nil {
		fmt.Println(fmt.Sprintf("Unable to post new trend entry for manufacturer %s, got error: %s", trendEntry.Manufacturer, requestErr.Error()))
		return requestErr
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Unable to post new trend entry for manufacturer %s, got status code: %d", trendEntry.Manufacturer, response.StatusCode))
		return errors.New(fmt.Sprintf("Unable to post trend entry for manufacturer %s", trendEntry.Manufacturer))
	}
	fmt.Println(fmt.Sprintf("Successfully posted new trend entry for manufacturer %s. Returning.", trendEntry.Manufacturer))
	return nil
}