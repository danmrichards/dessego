package middleware

import (
	"log"
	"net/http"
)

// LogRequest is a HTTP middleware that logs the incoming request to the given
// endpoint.
func LogRequest(endpoint string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%q request from client addr: %q", endpoint, r.RemoteAddr)

		h.ServeHTTP(w, r)
	}
}
