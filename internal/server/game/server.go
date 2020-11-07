package game

import (
	"fmt"
	"net"
	"net/http"

	"github.com/danmrichards/dessego/internal/crypto"
	"github.com/danmrichards/dessego/internal/transport"
)

// Server is a game server.
type Server struct {
	l  net.Listener
	r  *http.ServeMux
	h  *http.Server
	rd transport.RequestDecrypter
}

// NewServer returns a game server configured to run on the given host and port.
func NewServer(port string) (s *Server, err error) {
	s = &Server{
		r: http.NewServeMux(),
	}

	addr := net.JoinHostPort("", port)
	s.l, err = net.Listen("tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("net listen: %w", err)
	}

	s.rd, err = crypto.NewDecrypter(crypto.DefaultAESKey)
	if err != nil {
		return nil, err
	}

	s.routes()

	s.h = &http.Server{
		Addr:    addr,
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
