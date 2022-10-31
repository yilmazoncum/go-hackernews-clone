package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"main/hn"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"
)

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

		stories, err := getCachedStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

var cache []item
var cacheExpiration time.Time
var cacheMutex sync.Mutex

func getCachedStories(numStories int) ([]item, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	if time.Now().Sub(cacheExpiration) < 0 {
		fmt.Print("Cached Storyies Return !!!")
		return cache, nil
	}
	fmt.Print("getTopStories() ????")
	stories, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}
	cache = stories
	cacheExpiration = time.Now().Add(5 * time.Minute)
	return cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	ids, err := hn.GetTopItems()
	if err != nil {
		return nil, errors.New("Failed to load top stories")
	}

	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	type result struct {
		index int
		item  item
		err   error
	}

	resultCh := make(chan result)
	for i := 0; i < len(ids); i++ {
		go func(index int, id int) {
			hnItem, err := hn.GetItem(id)
			if err != nil {
				resultCh <- result{err: err}
			}
			resultCh <- result{index: index, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
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
