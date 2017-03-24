/**
	Go Multi Thread Web Parser
	Main file
*/

package main

import (
	"fmt"
	"time"
)

// Types block -------------------------------------------------------------------------
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}


type fetchStruct struct {
	title string
	body string
	filename string
	subdir string
	urls []string
}

// Fetcher is Fetcher that returns canned results.
type fetcher map[string]*fetchStruct


// Functions block ---------------------------------------------------------------------
func Crawl(url string, fetcher Fetcher) {
	body, urls, err := fetcher.Fetch(url)
	
	return
}

func (f fetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

func main() {
	// start HTTP server
	start := time.Now()
	fmt.Println("Starting Web server...")
	simple_http("65123")
	fmt.Println("Web server started ", time.Since(start))
	
	// fill structures
	
	// start crawl
	Crawl("http://golang.org/", fetcher)
	
	for {
		
	}
}


/**************************************************************************************
	Supplementary functions
*/
func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}