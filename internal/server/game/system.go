package game

import (
	"bytes"
	"net/http"

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
		data := new(bytes.Buffer)
		data.Write([]byte{0x01, 0x02})

		for _, m := range s.gs.Motd() {
			data.WriteString(m)
			data.WriteByte(0x00)
		}

		if err := transport.WriteResponse(w, 0x02, data.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) timeMsgHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// First byte
		// 0x00 - nothing
		// 0x01 - undergoing maintenance
		// 0x02 - online service has been terminated
		data := []byte{0x00, 0x00, 0x00}

		if err := transport.WriteResponse(w, 0x22, data); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
