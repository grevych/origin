// Description: This file exposes the public HTTP service for origin.
package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grevych/gobox/pkg/events"
	"github.com/grevych/gobox/pkg/log"
)

// service is a type that contains receiver functions that serve as handlers
// for individual HTTP routes.
type service struct{}

// Handler returns the main http handler for this service
func Handler() http.Handler {
	// svc := &service{}

	routes := mux.NewRouter()
	// routes.Handle("/ping", handlers.Endpoint("ping", svc.ping))
	// routes.Handle("/pong", handlers.Endpoint("pong", svc.pong))

	return routes
}

// Place any http handler functions for your service here
func (s service) ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Debug(r.Context(), "io write error", events.NewErrorInfo(err))
	}
}

func (s service) pong(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ping")); err != nil {
		log.Debug(r.Context(), "io write error", events.NewErrorInfo(err))
	}
}
