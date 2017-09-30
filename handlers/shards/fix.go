package shards

import (
	"net/http"
	"path"
	"strconv"

	"github.com/dabbotorg/panel/handlers/utils"
)

// Fix removes a fixed shard from the pool
func (handler *Handler) Fix(w http.ResponseWriter, r *http.Request) {
	arg := path.Base(r.URL.Path)

	id, err := strconv.Atoi(arg)
	if err != nil {
		http.Error(w, "id must be an int", http.StatusBadRequest)
		return
	}

	sesion, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	user := sesion.Values["id"].(string)

	ok := false
	for i, shard := range handler.ShardStore {
		if shard.ID == id {
			go handler.sendFixedWebhook(shard, user)
			handler.ShardStore = append(handler.ShardStore[:i], handler.ShardStore[i+1:]...)
			ok = true
		}
	}
	if ok == false {
		http.Error(w, "no shard found given that id", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/shards/list", http.StatusTemporaryRedirect)
}
