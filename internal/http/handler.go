// Description: This file exposes the public HTTP service for origin.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grevych/gobox/pkg/events"
	"github.com/grevych/gobox/pkg/log"
	"github.com/grevych/origin/pkg/httpx"
)

// service is a type that contains receiver functions that serve as handlers
// for individual HTTP routes.
type service struct{}

// Handler returns the main http handler for this service
func Handler() http.Handler {
	svc := &service{}

	routes := mux.NewRouter()
	routes.Handle("/headers", httpx.Endpoint("headers", svc.getHeaders)).Methods(http.MethodGet)
	routes.Handle("/headers/{headerID}", httpx.Endpoint("header", svc.getHeader)).Methods(http.MethodGet)
	return routes
}

// Place any http handler functions for your service here
func (s service) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Debug(r.Context(), "io write error", events.NewErrorInfo(err))
	}
}

func (s service) pong(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ping")); err != nil {
		log.Debug(r.Context(), "io write error", events.NewErrorInfo(err))
	}
}

func (s service) getHeaders(w http.ResponseWriter, r *http.Request) {
	serializedHeaders, err := json.Marshal(r.Header)
	if err != nil {
		log.Error(r.Context(), "get headers", events.NewErrorInfo(err))
		w.WriteHeader(http.StatusInternalServerError)
		// return error
		return
	}

	if _, err := w.Write([]byte(serializedHeaders)); err != nil {
		log.Error(r.Context(), "get headers", events.NewErrorInfo(err))
	}
}

func (s service) getHeader(w http.ResponseWriter, r *http.Request) {
	headerID, ok := mux.Vars(r)["headerID"]
	if !ok || headerID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	serializedHeader, err := json.Marshal(r.Header[headerID])
	if err != nil {
		log.Error(r.Context(), "get header", events.NewErrorInfo(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(serializedHeader)); err != nil {
		log.Error(r.Context(), "get header", events.NewErrorInfo(err))
	}
}
