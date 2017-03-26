/**
	Go Multi Thread Web Parser. 
	http server file
*/

package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"fmt"
)

var pwd, _ = os.Getwd()

type Page struct {
	Title string
	Body  template.HTML
}

func loadPage(title string) (*Page, error) {
	filename := pwd + "/web/" + title + ".html"
	body, err := ioutil.ReadFile(filename)
	bodys := template.HTML(body)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: bodys}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	renderTemplate(w, "view", p)
}

var templates = template.Must(template.ParseFiles(pwd+"/tmpl/view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, r.URL.Path[1:])
	}
}

func simple_http(port string) {
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        makeHandler(viewHandler),
		//ReadTimeout:    10 * time.Second,
		//WriteTimeout:   10 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}
	
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Printf("HTTP Server failed: ", err.Error())
		}
	}()
}