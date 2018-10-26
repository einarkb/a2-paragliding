package paragliding

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	id int
}

func URLHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func (server *Server) Start() {
	http.HandleFunc("/", URLHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
