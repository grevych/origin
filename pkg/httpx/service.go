package httpx

import (
	"context"
	"net/http"
)

type Service struct {
	App http.Handler
}

func (s *Service) Run(ctx context.Context, addr string) error {
	return http.ListenAndServe(addr, s.App)
}
