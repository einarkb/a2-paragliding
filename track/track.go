package track

import (
	"encoding/json"
	"io"
	"net/http"

	db "github.com/einarkb/paragliding/database"
	igc "github.com/marni/goigc"
)

type TrackMgr struct {
	DB *db.DB
}

// HandlerPostTrack is the handler for POST /api/track. it registers the track and replies with the id
func (tMgr *TrackMgr) HandlerPostTrack(w http.ResponseWriter, r *http.Request) {
	var postData map[string]string
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err == nil {
		_, err2 := igc.ParseLocation(postData["url"])
		if err2 != nil {
			http.Error(w, "could not get a track from url: "+postData["url"], http.StatusNotFound)
			return
		}
		id, added := tMgr.DB.Insert("tracks", postData["url"])
		if added {
			w.Header().Add("content-type", "application/json")
			json.NewEncoder(w).Encode(struct {
				ID string `json:"id"`
			}{id})
		} else {
			http.Error(w, "track already exists with id: "+id, http.StatusBadRequest)
		}
	} else if err == io.EOF {
		http.Error(w, "POST body is empty", http.StatusBadRequest)
	} else {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

// HandlerGetAllTracks is the handler for GET /api/track. it replies with an array of all track ids
func (tMgr *TrackMgr) HandlerGetAllTracks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

}
