package main

/* type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
} */

import (
	"fmt"
	"main/hn"
)

const apiBase = "https://hacker-news.firebaseio.com/v0"

func main() {

	/* ids, _ := hn.GetTopItems(apiBase)
	fmt.Print(ids) */
	item, _ := hn.GetItem(apiBase, 33387890)
	fmt.Print(item)

	/* tmpl := template.Must(template.ParseFiles("layout.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}

		tmpl.Execute(w, data)
	})

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	} */
}
