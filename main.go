package main

import (
	/*"fmt"
	"log"
	executor "github.com/andream16/yaggts/exec"*/
	parser "github.com/andream16/yaggts/parse"
)

func main() {
	/*fmt.Println("Starting yaggts . . .")
	results := func () bool {
		return executor.DownloadCSV("samsung galaxy")
	}(); if !results {
		log.Fatal("Unable to download csv.")
	}
	fmt.Println("Got the csv.")*/
	parser.GetTrend()
}
