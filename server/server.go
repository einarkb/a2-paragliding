package server

import (
	db "github.com/einarkb/paragliding/database"
)

type Server struct {
	db *db.DB
}

func (server *Server) Start() {
	server.db = &db.DB{URI: "mongodb://test:test12@ds141783.mlab.com:41783", Name: "a2-trackdb"}
	server.db.Connect()
}

func Hello() string {
	return "hello!"
}
