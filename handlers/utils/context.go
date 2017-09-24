package utils

import (
	"github.com/dabbotorg/panel/config"
	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

// AddAuthContext adds the context vars for an authorized user to a pongo2 context
func AddAuthContext(session *sessions.Session, ctx pongo2.Context, conf config.Config) {
	if session.Values["authed"] == "true" {
		ctx["authed"] = true
		ctx["username"] = session.Values["username"]
		ctx["discrim"] = session.Values["discrim"]
		ctx["id"] = session.Values["id"]
		ctx["admin"] = isAdmin(session.Values["id"], conf)
	}
}

func isAdmin(id interface{}, conf config.Config) bool {
	if id == "" {
		return false
	}

	var permission *config.Permission
	for _, perm := range conf.Permissions {
		if id == perm.ID {
			permission = &perm
		}
	}
	if permission == nil {
		return false
	} else if permission.Level < 10 {
		return false
	} else {
		return true
	}
}
