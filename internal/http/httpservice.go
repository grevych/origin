// Description: This file exposes the private HTTP service for origin
package http

import (
	"context"
)

// PrivateHTTPDependencies is used to inject dependencies into the HTTPService service
// activity. Great examples of integrations to be placed into here would be a database
// connection or perhaps a redis client that the service activity needs to use.
type PrivateHTTPDependencies struct {
}

// PrivateHTTPService handles internal http requests, suchs as metrics, health
// and readiness checks. This is required for ALL services to have.
type PrivateHTTPService struct {
	// handlers.Service

	listenHost string
	listenPort int
	deps       *PrivateHTTPDependencies
}

type Config interface {
	ListenHost() string
	PrivateHTTPPort() int
	PublicHTTPPort() int
}

// NewPrivateHTTPService creates a new HTTPService service activity.
func NewPrivateHTTPService(cfg Config, deps *PrivateHTTPDependencies) *PrivateHTTPService {
	return &PrivateHTTPService{
		listenHost: cfg.ListenHost(),
		listenPort: cfg.PrivateHTTPPort(),
		deps:       deps,
	}
}

// Run is the entrypoint for the HTTPService serviceActivity.
func (s *PrivateHTTPService) Run(ctx context.Context) error {
	// create a http handler (handlers.Service does metrics, health etc)
	//s.App = http.NotFoundHandler()
	// return s.Service.Run(ctx, fmt.Sprintf("%s:%d", s.cfg.ListenHost, s.cfg.HTTPPort))
	return nil
}

// PublicHTTPDependencies is used to inject dependencies into the PublicHTTPService
// service activity. Great examples of integrations to be placed into here would be
// a database connection or perhaps a redis client that the service activity needs to
// use.
type PublicHTTPDependencies struct {
}

// PublicHTTPService handles public http service calls.
type PublicHTTPService struct {
	// handlers.PublicService

	listenHost string
	listenPort int
	deps       *PublicHTTPDependencies
}

// NewPublicHTTPService creates a new PublicHTTPService service activity.
func NewPublicHTTPService(cfg Config, deps *PublicHTTPDependencies) *PublicHTTPService {
	return &PublicHTTPService{
		listenHost: cfg.ListenHost(),
		listenPort: cfg.PrivateHTTPPort(),
		deps:       deps,
	}
}

// Run starts the HTTP service at the host/port specified in the config.
func (s *PublicHTTPService) Run(ctx context.Context) error {
	// s.App = Handler()
	// return s.PublicService.Run(ctx, fmt.Sprintf("%s:%d", s.cfg.ListenHost, s.cfg.PublicHTTPPort))
	return nil
}
