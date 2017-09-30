package shards

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/flosch/pongo2"
)

// ViewReport renders the view to report a shard as broken
func (handler *Handler) ViewReport(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{}
	ctx["bots"] = handler.Meta.Bots

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)

	err = handler.Templates.Report.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
	}
}

// Report handles a request from the report form
func (handler *Handler) Report(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	bot := r.Form.Get("Bot")
	numberStr := r.Form.Get("Number")
	now := time.Now()

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "number must be an int", http.StatusBadRequest)
		return
	}

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
	}
	user := session.Values["id"].(string)

	shard := Shard{
		ID:        handler.ShardCounter,
		Bot:       bot,
		Number:    number,
		Timestamp: now,
		UserID:    user,
	}
	handler.ShardStore = append(handler.ShardStore, shard)
	handler.ShardCounter++

	go handler.sendBrokenWebhook(shard)

	http.Redirect(w, r, "/shards/list", http.StatusTemporaryRedirect)
}
