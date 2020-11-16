package game

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/danmrichards/dessego/internal/transport"
)

// Server is a gamestate server.
type Server struct {
	nl net.Listener
	r  *http.ServeMux
	h  *http.Server

	l zerolog.Logger

	rd transport.RequestDecrypter

	ps Players
	gs State
	ms Messages
}

// NewServer returns a gamestate server configured to run on the given host and port.
func NewServer(port string, rd transport.RequestDecrypter, ps Players, gs State, ms Messages, l zerolog.Logger) (s *Server, err error) {
	s = &Server{
		r:  http.NewServeMux(),
		rd: rd,
		ps: ps,
		gs: gs,
		ms: ms,
		l:  l,
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

// Serve accepts incoming gamestate connections.
func (s *Server) Serve() error {
	return s.h.Serve(s.nl)
}

// Close closes the gamestate server.
func (s *Server) Close() error {
	return s.h.Close()
}
