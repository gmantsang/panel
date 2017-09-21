package utils

import (
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

// AddAuthContext adds the context vars for an authorized user to a pongo2 context
func AddAuthContext(session *sessions.Session, ctx pongo2.Context) {
	if session.Values["authed"] == "true" {
		ctx["authed"] = true
		ctx["username"] = session.Values["username"]
		ctx["discrim"] = session.Values["discrim"]
	}
}
