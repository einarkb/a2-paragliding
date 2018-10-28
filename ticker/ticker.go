package ticker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/einarkb/paragliding/database"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// MgrTicker is the manager for the ticker part of things
type MgrTicker struct {
	DB      *db.DB
	PageCap int
}

// Response represents the ticker response
type Response struct {
	TLatest    int64               `json:"t_latest"`
	TStart     int64               `json:"t_start"`
	TStop      int64               `json:"t_stop"`
	TrackIDs   []objectid.ObjectID `json:"tracks"`
	Processing int64               `json:"processing"`
}

// HandlerLatestTick is the handler for "GET /api/ticker/latest"
// it responds with the timestamp of teh lastest added track
func (mgrTicker *MgrTicker) HandlerLatestTick(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	tracks, err := mgrTicker.DB.GetAllTracks()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(tracks) == 0 {
		fmt.Fprint(w, "No tracks")
		return
	}
	fmt.Fprint(w, tracks[len(tracks)-1].Timestamp)
}

// HandlerTicker is the handler for GET /api/ticker/
func (mgrTicker *MgrTicker) HandlerTicker(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	startTime := time.Now()
	tracks, err := mgrTicker.DB.GetAllTracks()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if len(tracks) == 0 {
		fmt.Fprint(w, "No tracks")
		return
	}
	nTracks := len(tracks)
	tickerResp := Response{}
	tickerResp.TLatest = tracks[nTracks-1].Timestamp
	tickerResp.TStart = tracks[0].Timestamp
	stopIndex := (mgrTicker.PageCap - 1)
	if stopIndex > nTracks-1 {
		stopIndex = nTracks - 1
	}
	if stopIndex < 0 {
		http.Error(w, "PageCap variable is not configured to a positive number", http.StatusInternalServerError)
		return
	}
	tickerResp.TStop = tracks[stopIndex].Timestamp
	for i := 0; i < mgrTicker.PageCap; i++ {
		tickerResp.TrackIDs = append(tickerResp.TrackIDs, tracks[i].ID)
	}
	tickerResp.Processing = int64(float64(time.Since(startTime)) / float64(time.Millisecond))
	json.NewEncoder(w).Encode(tickerResp)
}
