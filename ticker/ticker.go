package ticker

import (
	"fmt"
	"net/http"

	db "github.com/einarkb/paragliding/database"
)

// MgrTicker is the manager for the ticker part of things
type MgrTicker struct {
	DB *db.DB
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
	}
	fmt.Fprint(w, tracks[len(tracks)-1].Timestamp)
}
