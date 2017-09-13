package server

import "net/http"
import "github.com/gorilla/mux"

// Serve the panel
func Serve(c Config) error {

	router := mux.NewRouter()
	err := BuildRouter(router, c)
	if err != nil {
		return err
	}

	srv := http.Server{
		Addr:    c.Host,
		Handler: router,
	}

	return srv.ListenAndServe()
}
