package track

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson/objectid"

	db "github.com/einarkb/paragliding/database"
	igc "github.com/marni/goigc"
)

type TrackMgr struct {
	DB *db.DB
}

type TrackInfo struct {
	ID          objectid.ObjectID `bson:"_id" json:"-"`
	HDate       string            `bson:"H_date" json:"H_Date"`
	Pilot       string            `bson:"pilot" json:"pilot"`
	Glider      string            `bson:"glider" json:"glider"`
	GliderID    string            `bson:"glider_id" json:"glider_id"`
	TrackLength string            `bson:"track_length" json:"track_length"`
	TrackURL    string            `bson:"track_url" json:"track_url"`
}

// HandlerPostTrack is the handler for POST /api/track. it registers the track and replies with the id
func (tMgr *TrackMgr) HandlerPostTrack(w http.ResponseWriter, r *http.Request) {
	var postData map[string]string
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err == nil {
		track, err2 := igc.ParseLocation(postData["url"])
		if err2 != nil {
			http.Error(w, "could not get a track from url: "+postData["url"], http.StatusNotFound)
			return
		}
		trackInfo := TrackInfo{ID: objectid.New(), HDate: track.Date.String(), Pilot: track.Pilot,
			Glider: track.GliderType, GliderID: track.GliderID, TrackLength: CalculatedistanceFromPoints(track.Points),
			TrackURL: postData["url"]}
		id, added := tMgr.DB.Insert("tracks", trackInfo)
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

// CalculatedistanceFromPoints take a set of points and retunr the total distance
func CalculatedistanceFromPoints(points []igc.Point) string {
	d := 0.0
	for i := 0; i < len(points)-1; i++ {
		d += points[i].Distance(points[i+1])
	}
	return strconv.FormatFloat(d, 'f', 2, 64)
}
