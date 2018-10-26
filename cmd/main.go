package main

import (
	"fmt"
	"net/http"

	"github.com/einarkb/paragliding/server"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, server.Hello(), "r.URL.Path")
	})

	http.ListenAndServe(":80", nil)
}
