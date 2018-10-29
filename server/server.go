package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/einarkb/paragliding/admin"
	db "github.com/einarkb/paragliding/database"
	"github.com/einarkb/paragliding/ticker"
	"github.com/einarkb/paragliding/track"
	"github.com/einarkb/paragliding/webhook"
)

type Server struct {
	db         *db.DB
	trackMgr   *track.TrackMgr
	webhookMgr *webhook.WebHookMgr
	adminMgr   *admin.AdminMgr
	mgrTicker  *ticker.MgrTicker
	startTime  time.Time

	//map request type (eg. GET/POST) that contains map of acceptable urls and the function to handle each url
	urlHandlers map[string]map[string]func(http.ResponseWriter, *http.Request)
}

// Start starts the server
func (server *Server) Start() {
	server.startTime = time.Now()
	server.db = &db.DB{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	server.db.Connect()
	server.mgrTicker = &ticker.MgrTicker{DB: server.db, PageCap: 5}
	server.webhookMgr = &webhook.WebHookMgr{DB: server.db, Ticker: server.mgrTicker}
	server.trackMgr = &track.TrackMgr{DB: server.db, WHMgr: server.webhookMgr}
	server.adminMgr = &admin.AdminMgr{DB: server.db}
	server.initHandlers()

	http.HandleFunc("/", server.urlHandler)
	http.ListenAndServe(":80", nil)
}

func (server *Server) initHandlers() {
	//intializing maps
	server.urlHandlers = make(map[string]map[string]func(http.ResponseWriter, *http.Request))
	server.urlHandlers["GET"] = make(map[string]func(http.ResponseWriter, *http.Request))
	server.urlHandlers["POST"] = make(map[string]func(http.ResponseWriter, *http.Request))
	server.urlHandlers["DELETE"] = make(map[string]func(http.ResponseWriter, *http.Request))

	// registering handlers
	server.urlHandlers["GET"]["^/paragliding$"] = func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "paragliding/api", http.StatusSeeOther)
	}

	server.urlHandlers["GET"]["^/paragliding/api$"] = func(w http.ResponseWriter, r *http.Request) {
		type MetaData struct {
			Uptime  string `json:"uptime"`
			Info    string `json:"info"`
			Version string `json:"version"`
		}

		w.Header().Add("content-type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		encoder.Encode(MetaData{server.calculateUptime(), "Service for Paragliding tracks.", "v1"})
	}

	server.urlHandlers["POST"]["^/paragliding/api/track$"] = server.trackMgr.HandlerPostTrack
	server.urlHandlers["GET"]["^/paragliding/api/track$"] = server.trackMgr.HandlerGetAllTracks
	server.urlHandlers["GET"]["^/paragliding/api/track/[a-zA-Z0-9]{1,100}$"] = server.trackMgr.HandlerGetTrackByID
	server.urlHandlers["GET"]["^/paragliding/api/track/[a-zA-Z0-9]{1,50}/[a-zA-Z0-9_.-]{1,50}$"] = server.trackMgr.HandlerGetTrackFieldByID

	server.urlHandlers["GET"]["^/paragliding/api/ticker/latest$"] = server.mgrTicker.HandlerLatestTick
	server.urlHandlers["GET"]["^/paragliding/api/ticker/$"] = server.mgrTicker.HandlerTicker
	server.urlHandlers["GET"]["^/paragliding/api/ticker/[0-9]{1,20}$"] = server.mgrTicker.HandlerTickerByTimestamp

	server.urlHandlers["POST"]["^/paragliding/api/webhook/new_track/$"] = server.webhookMgr.HandlerNewTrackWebHook

	server.urlHandlers["GET"]["^/paragliding/admin/api/tracks_count$"] = server.adminMgr.HandlerTrackCount
	server.urlHandlers["DELETE"]["^/paragliding/admin/api/tracks$"] = server.adminMgr.HandlerDeleteAllTracks

}

// urHandler is reponsible for routing the different requests to the correct handler
func (server *Server) urlHandler(w http.ResponseWriter, r *http.Request) {
	handlerMap, exists := server.urlHandlers[r.Method]
	if !exists { // if not a request type we will handle (not GET, POST or DELETE in this case)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//fmt.Fprint(w, r.URL.Path)
	for url, hFunc := range handlerMap {
		res, _ := regexp.MatchString(url, r.URL.Path)
		if res {
			//fmt.Fprint(w, "huuuu")
			hFunc(w, r)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (server *Server) calculateUptime() string {
	dur := time.Since(server.startTime)

	sec := int(dur.Seconds()) % 60
	min := int(dur.Minutes()) % 60
	hour := int(dur.Hours()) % 24
	day := int(dur.Hours()/24) % 7
	month := int(dur.Hours()/24/7/4.34524) % 12
	year := int(dur.Hours() / 24 / 365.25)

	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", year, month, day, hour, min, sec)
}
