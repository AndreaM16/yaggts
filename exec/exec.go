package exec

import (
	"net/url"
	"os/exec"
)

var argv = [2]string{"python", "/home/andream16/go/src/github.com/andream16/yaggts/yaggts-selenium/main.py"}

// Takes a query, calls a python script to download a csv from Google Trends.
// Returns true if all went ok, false otherwise.
func DownloadCSV(query string) bool {
	return func() bool {
		return callScraper(&query)
	}()
}

// Executes the python script.
func callScraper(query *string) bool {
	_, scraperError := exec.Command(argv[0], argv[1], escapeString(query)).Output()
	return scraperError == nil
}

// Escapes the query.
func escapeString(q *string) string {
	return func () string {
		return url.PathEscape(*q)
	}()
}
