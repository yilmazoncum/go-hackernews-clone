package main

import (
	"flag"
	"fmt"
	"log"
	"main/hn"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"
)

const apiBase = "https://hacker-news.firebaseio.com/v0"

func main() {
	var port, numStories int
	flag.IntVar(&port, "port", 8080, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	templ := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, templ))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, templ *template.Template) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ids, err := hn.GetTopItems(apiBase)
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}

		var stories []item
		for _, id := range ids {
			hnItem, err := hn.GetItem(apiBase, id)
			if err != nil {
				continue
			}
			item := parseHNItem(hnItem)
			//itemA := item{Item: hnItem, Host: hnItem.URL}

			stories = append(stories, item)
			if len(stories) >= numStories {
				break
			}

		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}

		err = templ.Execute(w, data)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})

}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
