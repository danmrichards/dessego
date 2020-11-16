package game

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/service/msg"

	"github.com/danmrichards/dessego/internal/transport"
)

// legacyMessageLimit is the limit of legacy messages to return via the API.
const legacyMessageLimit = 5

func (s *Server) getBloodMsgHandler() http.HandlerFunc {
	type getBloodMsgReq struct {
		BlockID     int    `form:"blockID"`
		ReplayNum   int    `form:"replayNum"`
		CharacterID string `form:"characterID"`
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

		var bmr getBloodMsgReq
		if err = transport.DecodeRequest(s.rd, b, &bmr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		msgs := make([]msg.BloodMsg, 0, 10)

		// Player own messages.
		pm, err := s.ms.Player(bmr.CharacterID, bmr.BlockID, bmr.ReplayNum)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining := bmr.ReplayNum - len(pm)
		msgs = append(msgs, pm...)

		// Other player messages.
		opm, err := s.ms.NonPlayer(bmr.CharacterID, bmr.BlockID, remaining)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining -= len(opm)
		msgs = append(msgs, opm...)

		// Legacy messages.
		if (len(pm)+len(opm)) < legacyMessageLimit && remaining > 0 {
			lm, err := s.ms.Legacy(bmr.BlockID, bmr.ReplayNum)
			if err != nil {
				s.l.Err(err).Msg("")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			msgs = append(msgs, lm...)
		}

		s.l.Debug().Msgf(
			"%d blood messages block: %d character: %q",
			len(msgs), bmr.BlockID, bmr.CharacterID,
		)

		// TODO: Serialize messages.

		fmt.Println(msgs)
	}
}
