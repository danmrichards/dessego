package bootstrap

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog"
)

// Server is a bootstrap server.
type Server struct {
	gsHost string
	gs     map[string]string

	nl net.Listener
	r  *http.ServeMux
	h  *http.Server

	l zerolog.Logger
}

// NewServer returns a bootstrap server configured to run on the given host and port.
//
// The server will provide data for a gamestate to bootstrap and talk to the configured gamestate servers.
func NewServer(port, gsHost string, gameServers map[string]string, l zerolog.Logger) (s *Server, err error) {
	s = &Server{
		gsHost: gsHost,
		r:      http.NewServeMux(),
		gs:     gameServers,
		l:      l,
	}

	addr := net.JoinHostPort("", port)
	s.nl, err = net.Listen("tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("net listen: %w", err)
	}

	s.routes()

	s.h = &http.Server{
		Addr:    addr,
		Handler: s.r,
	}

	return s, nil
}

// Serve accepts incoming bootstrap connections.
func (s *Server) Serve() error {
	return s.h.Serve(s.nl)
}

// Close closes the bootstrap server.
func (s *Server) Close() error {
	return s.h.Close()
}
