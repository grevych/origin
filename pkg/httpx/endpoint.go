package httpx

import (
	"context"
	"mime"
	"net/http"

	"github.com/grevych/gobox/pkg/events"
	"github.com/grevych/gobox/pkg/log"
)

func Endpoint(name string, handler http.HandlerFunc) http.Handler {
	log.Info(context.Background(), "registering endpoint", log.F{"name": name})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		evt := &events.HTTPRequest{}
		evt.FillFieldsFromRequest(r)
		// TODO: Init tracer
		// TODO: Handle default 404
		log.Info(ctx, "http request", evt)
		// TODO: Handle panics
		rec := &ResponseRecorder{
			ResponseWriter: w,
		}
		handler.ServeHTTP(rec, r)
		log.Debug(ctx, "endpoint call", log.F{
			"requestID":       evt.RequestID,
			"responseHeaders": rec.headers,
			"responseBody":    string(rec.body),
		})
		evt.FillResponseInfo(len(rec.body), rec.statusCode)
		log.Info(ctx, "http response", evt)
	})
}

// - with custom middleware
func EndpointWithCustomMiddlewares(name string, handler http.HandlerFunc, pre, post []http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})
}

// - with authentication
func EndpointWithAuthentication() {}

// - with rate limit
func EndpointWithRateLimit() {}

// - with proxy
func EndpointWithProxy() {}

// - with async after response
// - with body available after response
// Allows body to be read after response is sent
func EndpointWithAsyncTask() {}

// - with error handling accordingly to API
func EndpointWithErrorHandler() {}

// - with redirect
func EndpointWithRedirect() {}

func JSONEndpoint(name string, handler http.HandlerFunc) http.Handler {
	endpoint := Endpoint(name, handler)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		w.Header().Add("Content-Type", "application/json")
		endpoint.ServeHTTP(w, r)
	})
}
