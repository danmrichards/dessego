package game

import (
	"log"
	"net/http"
)

func (s *Server) routes() {
	s.r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("unhandled request (%s %s) from %q to %q", r.Method, r.URL.Path, r.RemoteAddr, s.l.Addr())
	})

	s.r.HandleFunc("/cgi-bin/login.spd", s.loginHandler())
}
