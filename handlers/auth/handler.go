package auth

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Handler handles endpoints for authentication
type Handler struct {
	Store        sessions.Store
	ClientID     string
	ClientSecret string
}

// BuildRouter builds a subrouter for this controller
func (handler *Handler) BuildRouter(router *mux.Router) {
	r := router.PathPrefix("/link").Subrouter()
	r.HandleFunc("/discord", handler.EndFlow).
		Queries("code", "")
	r.HandleFunc("/discord", handler.StartFlow)
	r.HandleFunc("/logout", handler.Logout)
}
