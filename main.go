package main

import (
	"fmt"
	"log"
	executor "github.com/andream16/yaggts/exec"
	parser "github.com/andream16/yaggts/parse"
	//"github.com/andream16/yaggts/util"
	"github.com/andream16/yaggts/util"
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
		trend, trendErr := parser.GetTrend();
		if trendErr == nil && len(trend) > 0 {
			fmt.Println("Got a valid trend.")
			for _, k := range trend {
				fmt.Println(k.Value, k.Date)
			}
			//util.DeleteCsv()
		} else {
			fmt.Println(trendErr.Error())
		}
	}
}
