package server

import (
	"net/http"

	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/panel/handlers"
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
	r.Handle("/link/discord", &handlers.LinkHandler{
		Store:        store,
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
	})
	r.Handle("/logout", &handlers.LogoutHandler{
		Store: store,
	})

	r.Handle("/radios/{action}", &handlers.AuthorizedMiddleware{
		MinLevel: 1,
		Store:    store,
		Config:   c,
		Handler: &handlers.RadioHandler{
			Store:        store,
			Config:       c,
			ListTemplate: templates.RadioList,
			EditTemplate: templates.RadioEdit,
		},
	})

	return nil
}
