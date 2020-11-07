package game

import (
	"log"
	"net/http"
	"strings"

	"github.com/danmrichards/dessego/internal/transport"
)

func (s *Server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First byte
		// 0x00 - present EULA, create account (not working)
		// 0x01 - present MOTD, can be multiple
		// 0x02 - "Your account is currently suspended."
		// 0x03 - "Your account has been banned."
		// 0x05 - undergoing maintenance
		// 0x06 - online service has been terminated
		// 0x07 - network play cannot be used with this version
		data := "\x01" + "\x02" + strings.Join(s.gs.Motd(), "\x00") + "\x00"

		if err := transport.WriteResponse(w, 0x02, data); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
