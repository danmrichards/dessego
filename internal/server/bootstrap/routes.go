package bootstrap

func (s *Server) routes() {
	s.r.HandleFunc("/", s.handleBootstrap())
}
