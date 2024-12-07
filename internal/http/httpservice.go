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
	*httpx.Server

	deps *PrivateHTTPDependencies
}

type Config interface {
	ListenHost() string
	PrivateHTTPPort() int
	PublicHTTPPort() int
}

// NewPrivateHTTPServer creates a new HTTPService service activity.
func NewPrivateHTTPServer(cfg Config, deps *PrivateHTTPDependencies) *PrivateHTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost(), cfg.PrivateHTTPPort())
	return &PrivateHTTPServer{
		Server: httpx.NewServer(addr, http.NotFoundHandler()),
		deps:   deps,
	}
}

// Run is the entrypoint for the HTTPService serviceActivity.
func (s *PrivateHTTPServer) Run(ctx context.Context) error {
	// create a http handler (handlers.Service does metrics, health etc)
	return s.Server.Run(ctx)
}

func (s *PrivateHTTPServer) Close(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// PublicHTTPDependencies is used to inject dependencies into the PublicHTTPService
// service activity. Great examples of integrations to be placed into here would be
// a database connection or perhaps a redis client that the service activity needs to
// use.
type PublicHTTPDependencies struct {
}

// PublicHTTPServer handles public http service calls.
type PublicHTTPServer struct {
	*httpx.Server

	deps *PublicHTTPDependencies
}

// NewPublicHTTPServer creates a new PublicHTTPService service activity.
func NewPublicHTTPServer(cfg Config, deps *PublicHTTPDependencies) *PublicHTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.ListenHost(), cfg.PublicHTTPPort())
	return &PublicHTTPServer{
		Server: httpx.NewServer(addr, Handler()),
		deps:   deps,
	}
}

// Run starts the HTTP service at the host/port specified in the config.
func (s *PublicHTTPServer) Run(ctx context.Context) error {
	log.Info(ctx, "starting public http server", log.F{"address": s.Server.Addr})
	return s.Server.Run(ctx)
}

func (s *PublicHTTPServer) Close(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
