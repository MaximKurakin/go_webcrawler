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
    "errors"
    "sync"
    "runtime"
    "strconv"
)

type safeMap struct {
    umap    map[string]bool
    mux     sync.Mutex
}
var urlFin = safeMap{umap: make(map[string]bool)}

var baseUrl *url.URL

// Types block -------------------------------------------------------------------------
type Fetcher interface {
    // Fetch returns the body of URL and
    // a slice of URLs found on that page.
    Fetch(url string) (err error)
    URLfind() (err error)
    saveResult() (err error)
}


type fetchResult struct {
    url         url.URL
    body        []byte
    urls        []string
}

// Fetcher is Fetcher that returns canned results.
//type fetchResult map[string]*fetchStruct


// Functions block ---------------------------------------------------------------------
func Crawl(url string, c chan string, wg *sync.WaitGroup) {
    var fetcher fetchResult
fmt.Println("Crawling begin")
    defer wg.Done();
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
    for _, u := range fetcher.urls {
        wg.Add(1)
        fmt.Println("Send: ", u)
        //c <- u
        go func () { c <- u }()
    }
fmt.Println("Sent")
    return
}

func (f *fetchResult) Fetch(url string) (error) {
    // check the URL
    u, err := chkURL(url)
    if err != nil { return err }
    f.url = u
    
    // get the URL
    response, err := http.Get(f.url.String());
    if err != nil { return err }
    if response.StatusCode != 200 { return errors.New(fmt.Sprintf("Error: %v on %s", response.StatusCode, f.url.String())) }
    
    // get body
    f.body, err = ioutil.ReadAll(response.Body)
    if err != nil { return err }
    defer response.Body.Close()
    
    // get URLs
    err = f.URLfind()
    if err != nil { return err }
    
    //fmt.Println(f.urls)
    
    return nil
}

func (f *fetchResult) makeFileName(urlstr string) (string) {
    var ff *url.URL
    var err error
    
    if (urlstr != "") {
        ff, err = url.Parse(urlstr)
        if err != nil { return "" }
    }else{ ff = &f.url }

    u_path := filepath.Clean(ff.Path)
    subdir, base_path := filepath.Split(u_path)

    if ("" == u_path || string(filepath.Separator) == u_path) {
        base_path = "index.html"
    }
    
    extension := filepath.Ext(base_path)
    if extension != "" {
        base_path = strings.TrimSuffix(base_path, filepath.Ext(base_path))
    }else{
        extension = ".html"
    }
    
    return stripchars(filepath.Clean("./" + ff.Host + subdir + "/" + base_path + extension), ":*?\"<>|")
}

func (f *fetchResult) saveResult() error {
    filename := f.makeFileName("")
    subdir, _ := filepath.Split(filename)
    os.MkdirAll(filepath.Clean(subdir), 0777);
    
    return ioutil.WriteFile(filename, f.body, 0600)
}

func (f *fetchResult) URLfind() error {
    var urlstr string
    var ok bool
    
    z := html.NewTokenizer(bytes.NewReader(f.body))
    
    LOOPME:
    for {
        tt := z.Next()

        switch {
        case tt == html.ErrorToken:
            // End of the document, we're done
            //fmt.Println(f.url.String(), " ", f.urls)
            return nil
        case tt == html.StartTagToken || tt == html.SelfClosingTagToken:
            t := z.Token()
            urlstr = ""
            
            // Check if the token is an <a> tag
            switch t.Data {
                case "a", "link":
                    // Extract the href value, if there is one
                    ok, urlstr = getLink(t, "href")
                    if !ok { continue LOOPME }
                case "script", "img":
                    ok, urlstr = getLink(t, "src")
                    if !ok { continue LOOPME }
                default:
                    continue LOOPME
            }
            
            // check URL
            ur, err := chkURL(urlstr)
            if err != nil { continue LOOPME }
            
            // Replace URL in BODY with local one
            f.body = bytes.Replace(f.body, []byte("href=\""+urlstr+"\""), []byte("href=\""+makeLocalURL(urlstr)+"\""), -1)
            
            // make ABS url
            if ur.IsAbs() == true {
                urlstr = urlstr
            } else if ur.IsAbs() == false {
                urlstr = baseUrl.ResolveReference(&ur).String()
            } else if strings.HasPrefix(urlstr, "//") {
                ur.Scheme = baseUrl.Scheme
                ur.Host = baseUrl.Host
                urlstr = ur.String()
            } else if strings.HasPrefix(urlstr, "/") {
                ur.Scheme = baseUrl.Scheme
                ur.Host = baseUrl.Host
                urlstr = ur.String()
            } else {
                urlstr = urlstr
            }

            // make absolute URL
            /*if (!ur.IsAbs()) {
                ur.Scheme = f.url.Scheme
                ur.Host = f.url.Host
                urlstr = ur.String()
            }*/

            // check if exists and add to the channel
            wasIn := addUrlFin(&urlFin, urlstr)
            if !wasIn { f.urls = append(f.urls, urlstr) }
        }
    }

    return nil
}

// TODO: figire out this issue later
func makeLocalURL(oUrl string) (lUrl string) {
    lUrl = strings.TrimSuffix(oUrl, "/") + ".html"
    
    return lUrl
}

func addUrlFin(uarr *safeMap, ustr string) (ok bool) {
    uarr.mux.Lock()
    defer uarr.mux.Unlock()
    
    _, ok = uarr.umap[ustr]
    if !ok { uarr.umap[ustr] = true }
    
    return ok
}

func chkURL(ustr string) (url.URL, error) {
    u, err := url.Parse(ustr)
    if err != nil {    return url.URL{}, err}
    if len(u.Opaque) > 0 { return url.URL{}, errors.New("Opaque found in URL") }
    if u.IsAbs() == true {
        if (baseUrl.Scheme != u.Scheme || baseUrl.Host != u.Host) { return url.URL{}, nil }
    }
    
    return *u, nil
}

func getLink(t html.Token, attrName string) (ok bool, href string) {
    // Iterate over all of the Token's attributes until we find an "href"
    FMARK:
    for _, a := range t.Attr {
        if a.Key == attrName {
            href = a.Val
            ok = true
            break FMARK
        }
    }
    
    return ok, href
}


func main() {
    //bSize := 10 // default
    
    if len(os.Args) < 2 { fmt.Println("Please, specify web site URL"); os.Exit(2) }
    u := os.Args[1]
    
    //set how many processes (threads to use)
    if len(os.Args) > 2 {
        threadsNum, _ := strconv.Atoi(os.Args[2])
        runtime.GOMAXPROCS(threadsNum)
        //bSize = threadsNum
    }
    
    // start HTTP server
    start := time.Now()
    fmt.Println("Starting Web server...")
    simple_http("65123")
    fmt.Println("Web server started ", time.Since(start))
    
    // start Crawl
    start = time.Now()
    fmt.Println("Start Crawling of ", u)
    // define Schema, Host for given URL
    baseUrl, _ = url.Parse(u)
    
    // Initiate Chanel
    c := make(chan string)
    wg := &sync.WaitGroup{}
    // Push first run
    urlFin.umap[u] = true
    wg.Add(1)
    go Crawl(u, c, wg)
    
    // closer
    go func() { wg.Wait(); close(c) }()
    
    // receiver
    counter := 1
    for i := range c {
        counter++
        go Crawl(i, c, wg)
    }
    
    fmt.Println("Crawl finished... ", time.Since(start), " Crawled: ", counter)
    //panic("aaaa")
    var input string
    fmt.Scanln(&input)
}


func stripchars(str, chr string) string {
    return strings.Map(func(r rune) rune {
        if strings.IndexRune(chr, r) < 0 { return r }
        return -1
    }, str)
}