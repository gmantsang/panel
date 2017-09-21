package auth

import (
	"net/http"

	"github.com/dabbotorg/panel/handlers/utils"
)

// Logout handles /link/logout
func (handler *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := handler.Store.Get(r, utils.SessionName)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}
	session.Values["authed"] = nil
	session.Values["auth_token"] = nil
	session.Values["refresh_token"] = nil

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, utils.SessionFailed(err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
