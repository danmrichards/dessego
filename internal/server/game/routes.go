package game

import (
	"log"
	"net/http"
)

const routePrefix = "/cgi-bin"

func (s *Server) routes() {
	s.r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("unhandled request (%s %s) from %q to %q", r.Method, r.URL.Path, r.RemoteAddr, s.l.Addr())
	})

	s.r.HandleFunc(routePrefix+"/login.spd", s.loginHandler())
	s.r.HandleFunc(routePrefix+"/initializeCharacter.spd", s.initCharacterHandler())
}
