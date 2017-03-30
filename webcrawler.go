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
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"golang.org/x/net/html"
	"strings"
)

// Types block -------------------------------------------------------------------------
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (err error)
	URLfind() (err error)
	saveResult() (err error)
}


type fetchResult struct {
	url 		url.URL
	body		[]byte
	urls		[]string
	filePath	[]string
	subdir		string
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
	f.body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	
	// get URLs
	err = f.URLfind()
	if err != nil {
		return err
	}
	
	// get URL path
	f.url = *response.Request.URL
	
	return nil
}

func (f *fetchResult) makeFileName(urlstr string) (fileName string) {
	var ff *url.URL
	var err error
	
	if (urlstr != "") {
		ff, err = url.Parse(urlstr)
		if err != nil {
			return ""
		}
	}else{
		ff = &f.url
	}

	u_path := filepath.Clean(ff.Path)
	subdir, base_path := filepath.Split(u_path)

	if ("" == u_path || string(filepath.Separator) == u_path) {
		base_path = "index"
	}
	fileName = stripchars(filepath.Clean("./" + ff.Host + subdir + "/" + base_path + ".html"), ":*?\"<>|")

	return fileName
}

func (f *fetchResult) saveResult() error {
	filename := f.makeFileName("")
	subdir, _ := filepath.Split(filename)
	os.MkdirAll(filepath.Clean(subdir), 0777);
	
	return ioutil.WriteFile(filename, f.body, 0600)
}

func (f *fetchResult) URLfind() error {
	z := html.NewTokenizer(bytes.NewReader(f.body))
	
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			fmt.Println(f.filePath)
			return nil
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, urlstr := getHref(t)
			if !ok {
				continue
			}
			
			//// Make sure the url begines in http**
			//hasProto := strings.Index(url, "http") == 0
			//if hasProto {
			u, err := url.Parse(urlstr)
			if (!url.IsAbs(u)) {
				u.Scheme = f.url.Scheme
				u.Host = f.url.Host
			}
				f.urls = append(f.urls, urlstr)
				f.filePath = append(f.filePath, f.makeFileName(urlstr))
			//}
		}
	}

	return nil
}

func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	
	return ok, href
}

func main() {
	// start HTTP server
	start := time.Now()
	fmt.Println("Starting Web server...")
	simple_http("65123")
	fmt.Println("Web server started ", time.Since(start))
	
	// start crawl
	Crawl("http://localhost:65123/b")
	
	for {
		
	}
}


func stripchars(str, chr string) string {
    return strings.Map(func(r rune) rune {
        if strings.IndexRune(chr, r) < 0 {
            return r
        }
        return -1
    }, str)
}