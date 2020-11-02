package bootstrap

import (
	"fmt"
	"net"
	"net/http"
)

// Server is a bootstrap server.
type Server struct {
	gsHost string
	gs     map[string]string
	l      net.Listener
	r      *http.ServeMux
	h      *http.Server
}

// NewServer returns a bootstrap server configured to run on the given host and port.
//
// The server will provide data for a game to bootstrap and talk to the configured game servers.
func NewServer(port, gsHost string, gameServers map[string]string) (s *Server, err error) {
	s = &Server{
		gsHost: gsHost,
		r:      http.NewServeMux(),
		gs:     gameServers,
	}

	addr := net.JoinHostPort("", port)
	s.l, err = net.Listen("tcp4", addr)
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
	return s.h.Serve(s.l)
}

// Close closes the bootstrap server.
func (s *Server) Close() error {
	return s.h.Close()
}
