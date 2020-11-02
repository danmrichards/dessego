package game

import (
	"fmt"
	"net"
	"net/http"
)

// Server is a game server.
type Server struct {
	host string
	port string
	l    net.Listener
	r    *http.ServeMux
	h    *http.Server
}

// NewServer returns a game server configured to run on the given host and port.
func NewServer(host, port string) (s *Server, err error) {
	s = &Server{
		host: host,
		port: port,
		r:    http.NewServeMux(),
	}

	s.l, err = net.Listen("tcp4", net.JoinHostPort(host, port))
	if err != nil {
		return nil, fmt.Errorf("net listen: %w", err)
	}

	s.routes()

	s.h = &http.Server{
		Addr:    net.JoinHostPort(s.host, s.port),
		Handler: s.r,
	}

	return s, nil
}

// Serve accepts incoming game connections.
func (s *Server) Serve() error {
	return s.h.Serve(s.l)
}

// Close closes the game server.
func (s *Server) Close() error {
	return s.h.Close()
}
