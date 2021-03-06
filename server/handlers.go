package server

import (
	"log"
	"net/http"

	"github.com/dabbotorg/panel/handlers/shards"

	"github.com/dabbotorg/panel/config"
	"github.com/dabbotorg/panel/handlers"
	"github.com/dabbotorg/panel/handlers/auth"
	"github.com/dabbotorg/panel/handlers/radios"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// BuildRouter adds this website's routes
func BuildRouter(r *mux.Router, c config.Config, m config.Metadata) error {

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
		Config:   c,
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
		Meta:      m,
		Templates: radioTemplates,
	}
	radioHandler.BuildRouter(r)

	shardTemplates, err := shards.CompileTemplates()
	if err != nil {
		log.Fatalln(err)
	}
	shardHandler := &shards.Handler{
		Store:        store,
		Config:       c,
		Meta:         m,
		Templates:    shardTemplates,
		ShardStore:   make([]shards.Shard, 0),
		ShardCounter: 0,
	}
	shardHandler.BuildRouter(r)

	if c.Debug {
		r.HandleFunc("/compile", func(w http.ResponseWriter, r *http.Request) {
			radioTemplates, err = radios.CompileTemplates()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			radioHandler.Templates = radioTemplates
			shardTemplates, err := shards.CompileTemplates()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			shardHandler.Templates = shardTemplates
		})
	}

	return nil
}
