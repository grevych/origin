// Description: This file exposes the private HTTP service for origin
package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grevych/gobox/pkg/log"
	"github.com/grevych/origin/pkg/httpx"
)

// PrivateHTTPDependencies is used to inject dependencies into the HTTPService service
// activity. Great examples of integrations to be placed into here would be a database
// connection or perhaps a redis client that the service activity needs to use.
type PrivateHTTPDependencies struct {
}

// PrivateHTTPServer handles internal http requests, suchs as metrics, health
// and readiness checks. This is required for ALL services to have.
type PrivateHTTPServer struct {
	httpx.Service

	listenHost string
	listenPort int
	deps       *PrivateHTTPDependencies
}

type Config interface {
	ListenHost() string
	PrivateHTTPPort() int
	PublicHTTPPort() int
}

// NewPrivateHTTPServer creates a new HTTPService service activity.
func NewPrivateHTTPServer(cfg Config, deps *PrivateHTTPDependencies) *PrivateHTTPServer {
	return &PrivateHTTPServer{
		listenHost: cfg.ListenHost(),
		listenPort: cfg.PrivateHTTPPort(),
		deps:       deps,
	}
}

// Run is the entrypoint for the HTTPService serviceActivity.
func (s *PrivateHTTPServer) Run(ctx context.Context) error {
	// create a http handler (handlers.Service does metrics, health etc)
	s.App = http.NotFoundHandler()
	return s.Service.Run(ctx, fmt.Sprintf("%s:%d", s.listenHost, s.listenPort))
}

// PublicHTTPDependencies is used to inject dependencies into the PublicHTTPService
// service activity. Great examples of integrations to be placed into here would be
// a database connection or perhaps a redis client that the service activity needs to
// use.
type PublicHTTPDependencies struct {
}

// PublicHTTPServer handles public http service calls.
type PublicHTTPServer struct {
	httpx.Service

	listenHost string
	listenPort int
	deps       *PublicHTTPDependencies
}

// NewPublicHTTPServer creates a new PublicHTTPService service activity.
func NewPublicHTTPServer(cfg Config, deps *PublicHTTPDependencies) *PublicHTTPServer {
	return &PublicHTTPServer{
		listenHost: cfg.ListenHost(),
		listenPort: cfg.PublicHTTPPort(),
		deps:       deps,
	}
}

// Run starts the HTTP service at the host/port specified in the config.
func (s *PublicHTTPServer) Run(ctx context.Context) error {
	s.App = Handler()
	log.Info(ctx, "starting public http server", log.F{"host": s.listenHost, "port": s.listenPort})
	return s.Service.Run(ctx, fmt.Sprintf("%s:%d", s.listenHost, s.listenPort))
}
