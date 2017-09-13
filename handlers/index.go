package handlers

import (
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

// IndexHandler handles /
type IndexHandler struct {
	Store    sessions.Store
	Template *pongo2.Template
}

func (handler *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	session, err := handler.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
	}

	ctx := pongo2.Context{}
	if session.Values["authed"] == "true" {
		ctx["authed"] = true
		ctx["username"] = session.Values["username"]
		ctx["discrim"] = session.Values["discrim"]
	}

	err = handler.Template.ExecuteWriter(ctx, w)
	if err != nil {
		http.Error(w, TemplateFailed(err), http.StatusInternalServerError)
	}
}
