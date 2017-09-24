package radios

import (
	"net/http"

	"github.com/dabbotorg/panel/handlers/utils"

	"github.com/dabbotorg/panel/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Handler manages radio stations
type Handler struct {
	Store     sessions.Store
	Templates *Templates
	Config    config.Config
	Meta      config.Metadata
}

// BuildRouter builds a subrouter for this controller
func (handler *Handler) BuildRouter(router *mux.Router) {
	r := router.PathPrefix("/radios").Subrouter()
	r.HandleFunc("/list/{state:valid|escrow|broken}", handler.ViewList)

	r.HandleFunc("/edit/{name}", handler.ViewEdit).
		Methods(http.MethodGet)
	r.Handle("/edit/{name}", &utils.AuthorizedMiddleware{
		Config:      handler.Config,
		MinLevel:    5,
		Form:        true,
		HandlerFunc: handler.Edit,
	}).Methods(http.MethodPost)

	r.HandleFunc("/create/", handler.ViewCreate).
		Methods(http.MethodGet)
	r.Handle("/create/", &utils.AuthorizedMiddleware{
		Config:      handler.Config,
		MinLevel:    5,
		Form:        true,
		HandlerFunc: handler.Create,
	}).Methods(http.MethodPost)

	r.HandleFunc("/delete/{name}", handler.ViewDelete).
		Methods(http.MethodGet)
	r.Handle("/delete/{name}", &utils.AuthorizedMiddleware{
		Config:      handler.Config,
		MinLevel:    10,
		Form:        true,
		HandlerFunc: handler.Delete,
	}).Methods(http.MethodPost)

	r.HandleFunc("/break/{name}", handler.ViewBreak).
		Methods(http.MethodGet)
	r.Handle("/break/{name}", &utils.AuthorizedMiddleware{
		Config:      handler.Config,
		MinLevel:    10,
		Form:        true,
		HandlerFunc: handler.Break,
	}).Methods(http.MethodPost)

	r.Handle("/valid/{name}", &utils.AuthorizedMiddleware{
		Config:      handler.Config,
		MinLevel:    10,
		Form:        true,
		HandlerFunc: handler.Valid,
	})
}
