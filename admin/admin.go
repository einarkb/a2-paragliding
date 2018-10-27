package admin

import (
	"fmt"
	"net/http"

	db "github.com/einarkb/paragliding/database"
)

type AdminMgr struct {
	DB *db.DB
}

// HandlerTrackCount is the handler for GET /admin/api/tracks_count. it responds with the number of trascks in the database
func (aMgr *AdminMgr) HandlerTrackCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	fmt.Fprint(w, aMgr.DB.GetTrackCount())
}
