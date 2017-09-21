package server

import (
	"log"
	"net/http"

	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/panel/handlers"
	"github.com/dabbotorg/panel/handlers/auth"
	"github.com/dabbotorg/panel/handlers/radios"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// BuildRouter adds this website's routes
func BuildRouter(r *mux.Router, c config.Config) error {

	store := sessions.NewCookieStore([]byte(c.Secret))
	templates, err := CompileTemplates()
	if err != nil {
		return err
	}

	r.PathPrefix("/assets/").
		Handler(http.FileServer(http.Dir("/assets/")))

	r.Handle("/", &handlers.IndexHandler{
		Store:    store,
		Template: templates.Index,
	})

	authHandler := &auth.Handler{
		Store:        store,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	}
	authHandler.BuildRouter(r)

	radioTemplates, err := radios.CompileTemplates()
	if err != nil {
		log.Fatalln(err)
	}
	radioHandler := &radios.Handler{
		Store:     store,
		Config:    c,
		Templates: radioTemplates,
	}
	radioHandler.BuildRouter(r)

	return nil
}
