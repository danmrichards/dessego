package game

import (
	"log"
	"net/http"

	"github.com/danmrichards/dessego/internal/transport"
)

func (s *Server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("login request from %q to %q", r.RemoteAddr, s.l.Addr())

		// TODO: Extract MoTD to game state manager.

		// Message of the Day.
		motd := "Welcome to DeSSE Go\r\n"
		motd += "A server emulator for Demon's Souls implemented in Go\r\n"
		motd += "Source code:\r\n"
		motd += "https://github.com/danmrichards/dessego\r\n"

		motd2 := "TODO: Add server stats here"

		// first byte
		// 0x00 - present EULA, create account (not working)
		// 0x01 - present MOTD, can be multiple
		// 0x02 - "Your account is currently suspended."
		// 0x03 - "Your account has been banned."
		// 0x05 - undergoing maintenance
		// 0x06 - online service has been terminated
		// 0x07 - network play cannot be used with this version
		cmd := 0x02
		data := "\x01" + "\x02" + motd + "\x00" + motd2 + "\x00"

		if err := transport.WriteResponse(w, cmd, data); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
