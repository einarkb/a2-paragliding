package webhook

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	db "github.com/einarkb/paragliding/database"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// WebHookMgr is the manager for webhooks
type WebHookMgr struct {
	DB *db.DB
}

// HandlerNewTrackWebHook is the handler for POST /api/webhook/new_track/.
// it registers a new webhook and reponds with the id assigned to it
func (whMgr *WebHookMgr) HandlerNewTrackWebHook(w http.ResponseWriter, r *http.Request) {
	var postData map[string]string
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err == nil {
		triggerVal, err2 := strconv.Atoi(postData["minTriggerValue"])
		if err2 != nil {
			http.Error(w, "triggervalue is not a number", http.StatusBadRequest)
			return
		}
		wekbookInfo := db.WebhookInfo{ID: objectid.New(), WebhookURL: postData["url"], MinTriggerValue: triggerVal, Counter: 0}
		id, added := whMgr.DB.Insert("webhooks", wekbookInfo)
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

func InvokeWebHooks() {
	//todo timestamps and multiple ids

}
