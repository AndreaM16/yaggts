package main

import (
	"fmt"
	executor "github.com/andream16/yaggts/exec"
)

func main() {
	fmt.Println("Starting yaggts . . .")
	results := func () bool {
		return executor.DownloadCSV("samsung galaxy")
	}(); fmt.Println(results)
}
