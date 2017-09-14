package handlers

import (
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

func addAuthContext(session *sessions.Session, ctx pongo2.Context) {
	if session.Values["authed"] == "true" {
		ctx["authed"] = true
		ctx["username"] = session.Values["username"]
		ctx["discrim"] = session.Values["discrim"]
	}
}
