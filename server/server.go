package server

import (
	"net/http"

	"github.com/dabbotorg/panel/config"
	"github.com/gorilla/mux"
)

// Serve the panel
func Serve(c config.Config, m config.Metadata) error {

	router := mux.NewRouter()
	err := BuildRouter(router, c, m)
	if err != nil {
		return err
	}

	srv := http.Server{
		Addr:    c.Host,
		Handler: router,
	}

	return srv.ListenAndServe()
}
