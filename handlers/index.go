package handlers

import (
	"net/http"

	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

// IndexHandler handles /
type IndexHandler struct {
	Store    sessions.Store
	Template *pongo2.Template
}

func (handler *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
	}

	ctx := pongo2.Context{}
	utils.AddAuthContext(session, ctx)

	err = handler.Template.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, utils.TemplateFailed(err), http.StatusInternalServerError)
	}
}
