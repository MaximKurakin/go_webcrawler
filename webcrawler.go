/**
	Go Multi Thread Web Parser
	Main file
*/

package main

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"os"
)

// Types block -------------------------------------------------------------------------
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body []byte, urls []string, err error)
}


type fetchResult struct {
	url string
	body []byte
	urls []string
}

// Fetcher is Fetcher that returns canned results.
//type fetchResult map[string]*fetchStruct


// Functions block ---------------------------------------------------------------------
func Crawl(url string, fetcher Fetcher) {
return
	body, urls, err := fetcher.Fetch(url)
_ = urls
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("found: %s %q\n", url, body)
	var result = fetchResult{
		"/p",
		body,
		nil,
	}
	_ = result
	err = result.saveResult()
	if err != nil {
		fmt.Println(err)
		return
	}
	// recursion
	/*for _, u := range urls {
		Crawl(u, fetcher)
	}*/

	return
}

func (f fetchResult) Fetch(url string) ([]byte, []string, error) {
	
	response, err := http.Get(url);
	if err != nil {
		return nil, nil, err
	}
	
	// get body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	
	
	return body, nil, nil
}

func (p *fetchResult) saveResult() error {
	filename := "." + "/localhost65123" + p.url + ".html"
	os.MkdirAll("." + "/localhost65123" + "/", 0777);
	//fmt.Println(filename)
	return ioutil.WriteFile(filename, p.body, 0600)
}

func main() {
	// start HTTP server
	start := time.Now()
	fmt.Println("Starting Web server...")
	simple_http("65123")
	fmt.Println("Web server started ", time.Since(start))
	
	// fill structures
	var fetcher fetchResult
	// start crawl
	Crawl("http://localhost:65123/p", fetcher)
	
	for {
		
	}
}


/**************************************************************************************
	Supplementary functions
*/
func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}