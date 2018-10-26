package server

import (
	"fmt"
	"net/http"

	db "github.com/einarkb/paragliding/database"
)

type Server struct {
	db *db.DB

	//map request type (eg. GET/POST) that contains map of acceptable urls and the function to handle each url
	urlHandlers map[string]map[string]func(http.ResponseWriter, *http.Request)
}

// Start starts the server
func (server *Server) Start() {
	server.db = &db.DB{URI: "mongodb://test:test12@ds141783.mlab.com:41783", Name: "a2-trackdb"}
	server.db.Connect()
	server.initHandlers()

	http.HandleFunc("/", server.urlHandler)
	http.ListenAndServe(":80", nil)
}

func (server *Server) initHandlers() {
	server.urlHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	server.urlHandlers["GET"]["/test/"] = handleTest
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello!!!!")
}

func (server *Server) urlHandler(w http.ResponseWriter, r *http.Request) {
	/*handlerMap, exists := server.urlHandlers[r.Method]
	if !exists {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	for url, hFunc := range handlerMap {
		fmt.Fprint(w, r.URL.Path)
		if r.URL.Path == url {
			hFunc(w, r)
			return
		}
	}*/
	fmt.Fprint(w, r.URL.Path)

}
