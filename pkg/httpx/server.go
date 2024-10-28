package httpx

import (
	"context"
	"net/http"
)

type Server struct {
	*http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	return &Server{
		&http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	status := make(chan error, 1)
	go func() {
		status <- s.ListenAndServe()
	}()

	select {
	case err := <-status:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
