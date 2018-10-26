package main

import (
	"fmt"
	"net/http"

	s "github.com/einarkb/paragliding/server"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, s.hello(), "r.URL.Path")
	})

	http.ListenAndServe(":80", nil)
}
