package server

import (
	"fmt"
	"net/http"
	"time"

	db "github.com/einarkb/paragliding/database"
)

type Server struct {
	db         *db.DB
	apiInfoMgr *APIInfoMgr

	//map request type (eg. GET/POST) that contains map of acceptable urls and the function to handle each url
	urlHandlers map[string]map[string]func(http.ResponseWriter, *http.Request)
}

// Start starts the server
func (server *Server) Start() {
	server.db = &db.DB{URI: "mongodb://test:test12@ds141783.mlab.com:41783", Name: "a2-trackdb"}
	server.db.Connect()
	server.apiInfoMgr = &APIInfoMgr{startTime: time.Now()}
	server.initHandlers()

	http.HandleFunc("/", server.urlHandler)
	http.ListenAndServe(":80", nil)
}

func (server *Server) initHandlers() {
	server.urlHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	server.urlHandlers["GET"] = make(map[string]func(http.ResponseWriter, *http.Request))
	server.urlHandlers["POST"] = make(map[string]func(http.ResponseWriter, *http.Request))

	server.urlHandlers["GET"]["/test/"] = handleTest
	server.urlHandlers["GET"]["/paragliding"] = func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "paragliding/api", http.StatusSeeOther)
	}
	server.urlHandlers["GET"]["/paragliding/api"] = server.apiInfoMgr.APIHandler
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello!!!!")
}

// urHandler is reponsible for routing the different requests to the correct handler
func (server *Server) urlHandler(w http.ResponseWriter, r *http.Request) {
	handlerMap, exists := server.urlHandlers[r.Method]
	if !exists { // if not a request type we will handle (not GET or POST in this case)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	for url, hFunc := range handlerMap {
		if r.URL.Path == url {
			hFunc(w, r)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
