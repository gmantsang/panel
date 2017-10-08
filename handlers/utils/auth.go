package utils

import (
	"net/http"

	"github.com/dabbotorg/panel/config"
	"github.com/gorilla/sessions"
)

// AuthorizedMiddleware protects a route that requires authentication
type AuthorizedMiddleware struct {
	MinLevel    int
	Form        bool
	Store       sessions.Store
	Handler     *http.Handler
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
	Config      config.Config
}

func (middleware *AuthorizedMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := middleware.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}

	id := session.Values["id"]
	if id == "" {
		http.Error(w, NotAuthorized, http.StatusForbidden)
		return
	}

	var permission *config.Permission
	for _, perm := range middleware.Config.Permissions {
		if id == perm.ID {
			permission = &perm
			break
		}
	}
	if permission == nil {
		http.Error(w, NotAuthorized, http.StatusForbidden)
		return
	}
	if permission.Level < middleware.MinLevel {
		http.Error(w, NotAuthorized, http.StatusForbidden)
		return
	}

	if middleware.Handler != nil {
		(*middleware.Handler).ServeHTTP(w, r)
		return
	}
	middleware.HandlerFunc(w, r)
}
