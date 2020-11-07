package bootstrap

import "github.com/danmrichards/dessego/internal/server/middleware"

func (s *Server) routes() {
	s.r.HandleFunc("/", middleware.LogRequest("bootstrap", s.bootstrapHandler()))
}
