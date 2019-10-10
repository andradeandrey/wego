package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Printf("Start Crawler")

	cwl, _ := NewCrawler("http://localhost:9999/chamada/fake", 100)
	go cwl.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":9999", nil)

}
