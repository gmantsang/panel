package handlers

import (
	"fmt"
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
	err := handler.Template.ExecuteWriter(nil, w)
	if err != nil {
		http.Error(w, fmt.Sprintf(TemplateFailed, err), http.StatusInternalServerError)
	}
}
