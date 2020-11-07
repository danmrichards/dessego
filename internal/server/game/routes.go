package game

import (
	"log"
	"net/http"

	"github.com/danmrichards/dessego/internal/server/middleware"
)

const routePrefix = "/cgi-bin"

func (s *Server) routes() {
	s.r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("unhandled request (%s %s) from %q to %q", r.Method, r.URL.Path, r.RemoteAddr, s.l.Addr())
	})

	s.r.HandleFunc(
		routePrefix+"/login.spd",
		middleware.LogRequest("login", s.loginHandler()),
	)
	s.r.HandleFunc(
		routePrefix+"/initializeCharacter.spd",
		middleware.LogRequest("init character", s.initCharacterHandler()),
	)
}
