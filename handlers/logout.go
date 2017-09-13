package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// LogoutHandler handles logging out
type LogoutHandler struct {
	Store sessions.Store
}

func (handler *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := handler.Store.Get(r, SessionName)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}
	session.Values["authed"] = nil
	session.Values["auth_token"] = nil
	session.Values["refresh_token"] = nil

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, SessionFailed(err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
