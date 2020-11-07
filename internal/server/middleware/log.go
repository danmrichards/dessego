package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

// LogRequest is a HTTP middleware that logs the incoming request to the given
// endpoint.
func LogRequest(l zerolog.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("client", r.RemoteAddr).
			Msg("")

		h.ServeHTTP(w, r)
	}
}
