package httpx

import (
	"net/http"
)

// ResponseRecorder wraps a http.ResponseWriter to
// record its status code before the response is written.
// TODO: Explain why it does not matter to lose the rest of the
// default interfaces.
type ResponseRecorder struct {
	http.ResponseWriter

	statusCode int
	body       []byte
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
	}
}

func (rw *ResponseRecorder) Write(p []byte) (n int, err error) {
	rw.body = p
	return rw.ResponseWriter.Write(p)
}

func (rw *ResponseRecorder) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *ResponseRecorder) WriteHeader(statusCode int) {
	rw.statusCode = statusCode

	if statusCode == http.StatusForbidden {
		rw.ResponseWriter.WriteHeader(http.StatusNotFound)
		return
	}

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseRecorder) GetBody() []byte {
	return rw.body
}

func (rw *ResponseRecorder) GetStatusCode() int {
	return rw.statusCode
}
