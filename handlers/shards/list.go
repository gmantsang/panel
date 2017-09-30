package shards

import (
	"net/http"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/flosch/pongo2"
)

// ViewList renders the list of shards
func (handler *Handler) ViewList(w http.ResponseWriter, r *http.Request) {
	ctx := pongo2.Context{}

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	utils.AddAuthContext(session, ctx, handler.Config)

	ctx["shards"] = handler.ShardStore

	err = handler.Templates.List.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
	}
}
