package radios

import (
	"net/http"

	"github.com/dabbotorg/panel/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Handler manages radio stations
type Handler struct {
	Store     sessions.Store
	Templates *Templates
	Config    config.Config
}

// BuildRouter builds a subrouter for this controller
func (handler *Handler) BuildRouter(router *mux.Router) {
	r := router.PathPrefix("/radios").Subrouter()
	r.HandleFunc("/list/{state:valid|escrow|broken}", handler.ViewList)
	r.HandleFunc("/edit/{name}", handler.ViewEdit).
		Methods(http.MethodGet)
	r.HandleFunc("/edit/{name}", handler.Edit).
		Methods(http.MethodPost)
}
