package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const apiBase = "https://hacker-news.firebaseio.com/v0"

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`
}

func GetTopItems(apiBase string) ([]int, error) {
	response, err := http.Get(fmt.Sprintf("%s/topstories.json", apiBase))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var ids []int
	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func GetItem(apiBase string, id int) (Item, error) {
	var item Item

	response, err := http.Get(fmt.Sprintf("%s/item/%d.json", apiBase, id))
	if err != nil {
		return item, err
	}
	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	err = dec.Decode(&item)
	if err != nil {
		return item, err
	}

	return item, nil
}
