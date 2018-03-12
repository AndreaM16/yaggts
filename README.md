# yaggts
Yet Another Golang Google Trends Scraper

# How to run
1. Install python and go
2. Install selenium module with `sudo install pip selenium`
3. Download last geckodriver from https://github.com/mozilla/geckodriver/releases
4. Extract it and put it in /usr/bin
5. Change paths in `scraper/main.py` and in `exec/exec.go` for your system
6. Run with `go run main.go query`, for instance `go run main.go Apple`
