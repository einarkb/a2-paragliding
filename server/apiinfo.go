package server

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIInfoMgr stores the start time of the server/api. it had handler functions for api
type APIInfoMgr struct {
	startTime time.Time
}

// APIHandler is the handler for a get request to "/paragliding/api", it responds with the api's metadata
func (apiInfo *APIInfoMgr) APIHandler(w http.ResponseWriter, r *http.Request) {
	type MetaData struct {
		Uptime  string `json:"uptime"`
		Info    string `json:"info"`
		Version string `json:"version"`
	}

	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	encoder.Encode(MetaData{time.Since(apiInfo.startTime).String(), "Service for IGC tracks.", "v1"})
}
