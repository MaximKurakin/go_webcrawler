/**
	Go Multi Thread Web Parser
	Main file
*/

package main

import (
	"fmt"
	"time"
	"net/http"
	"net/url"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Types block -------------------------------------------------------------------------
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (err error)
	saveResult() (err error)
}


type fetchResult struct {
	url url.URL
	body []byte
	urls []string
}

// Fetcher is Fetcher that returns canned results.
//type fetchResult map[string]*fetchStruct


// Functions block ---------------------------------------------------------------------
func Crawl(url string) {
	var fetcher fetchResult
	
	// fetch URL
	err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// save file
	err = fetcher.saveResult()
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

func (f *fetchResult) Fetch(url string) (error) {
	// get the URL
	response, err := http.Get(url);
	if err != nil {
		return err
	}
	
	// get body
	// defer response.Body.Close() // seems don't need to explicitly call this method anymore
	
f.body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	
	// get URL path
	f.url = *response.Request.URL
	//fmt.Printf("%#v\n", f.url)
	
	return nil
}

func (p *fetchResult) saveResult() error {
	u_path		:= p.url.Path
	base_path	:= filepath.Base(u_path)
	subdir		:= filepath.Dir(u_path)
	
	if ("" == u_path || "/" == u_path) {
		base_path = "/index"
	}
	
	filename := filepath.Clean("./" + p.url.Host + subdir + "/" + base_path + ".html")
	os.MkdirAll(filepath.Clean("./" + p.url.Host + subdir), 0777);
	fmt.Println(filename, base_path, subdir)
	
	return ioutil.WriteFile(filename, p.body, 0600)
}

func main() {
	// start HTTP server
	start := time.Now()
	fmt.Println("Starting Web server...")
	simple_http("65123")
	fmt.Println("Web server started ", time.Since(start))
	
	// start crawl
	Crawl("http://ydacha-mb.org/arhiv-dokumentatsii-zastrojshhika")
	
	for {
		
	}
}


/**************************************************************************************
	Supplementary functions
*/
func remove(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}