package shards

import (
	"net/http"

	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/panel/handlers/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Handler manages shards
type Handler struct {
	Store        sessions.Store
	Templates    *Templates
	Config       config.Config
	Meta         config.Metadata
	ShardStore   []Shard
	ShardCounter int
}

// BuildRouter adds the route for shards
func (handler *Handler) BuildRouter(r *mux.Router) {
	r.HandleFunc("/shards/list", handler.ViewList)

	r.HandleFunc("/shards/report", handler.ViewReport).
		Methods(http.MethodGet)
	r.Handle("/shards/report", &utils.AuthorizedMiddleware{
		Store:       handler.Store,
		Config:      handler.Config,
		MinLevel:    5,
		HandlerFunc: handler.Report,
	})

	r.Handle("/shards/fix/{name}", &utils.AuthorizedMiddleware{
		Store:       handler.Store,
		Config:      handler.Config,
		MinLevel:    5,
		HandlerFunc: handler.Fix,
	})
}
