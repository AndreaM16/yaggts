package main

import (
	"flag"
	"fmt"
	"github.com/andream16/yaggts/util"
	"github.com/andream16/yaggts/model"
	executor "github.com/andream16/yaggts/exec"
	parser "github.com/andream16/yaggts/parse"
)

func main() {
	fmt.Println("Starting yaggts . . .")
	flags := parseFlags()
	manufacturers, manufacturersError := util.GetManufacturers(flags); if manufacturersError != nil {
		panic(manufacturersError)
	}
	for _, manufacturer := range manufacturers.Manufacturers {
		results := func() bool {
			return executor.DownloadCSV(manufacturer.Manufacturer)
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
				var dbEntry model.DbEntry
				trendEntry.Manufacturer = manufacturer.Manufacturer
				dbEntry.Date = k.Date
				dbEntry.Value = k.Value
				trendEntry.Trend = append(trendEntry.Trend, dbEntry)
			}
			util.PostTrend(trendEntry)
			util.DeleteCsv()
		} else {
			fmt.Println(trendErr.Error())
		}
	}
}

func parseFlags() model.Flag {
	var flags model.Flag
	flags.Host = flag.String("host", "localhost", "host")
	flags.Port = flag.Int("port", 8080, "port")
	flags.Route = flag.String("route", "manufacturer", "route")
	flags.Page = flag.Int("page", 1, "page")
	flags.Size = flag.Int("size", 100, "size")
	flag.Parse()
	return flags
}