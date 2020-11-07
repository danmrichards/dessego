package game

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/transport"
)

func (s *Server) initCharacterHandler() http.HandlerFunc {
	type initCharacterReq struct {
		CharacterID string `form:"characterID"`
		Index       int    `form:"index"`
		Version     int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var icr initCharacterReq
		if err = transport.DecodeRequest(s.rd, b, &icr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create the player, if it does not exist, in the DB.
		if err = s.ps.EnsureCreate(icr.CharacterID, icr.Index); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unique character ID.
		ucID := fmt.Sprintf("%s%d", icr.CharacterID, icr.Index)

		// Track the player in game state.
		s.gs.AddPlayer(r.RemoteAddr, ucID)

		cmd := 0x17
		data := ucID + "\x00"

		if err = transport.WriteResponse(w, cmd, data); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
