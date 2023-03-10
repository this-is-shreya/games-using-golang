package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "listening at %s", r.URL.Path)
		// fmt.Println("running")
	})

	http.ListenAndServe(":8000", nil)
}
