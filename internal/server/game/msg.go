package game

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/service/msg"
	"github.com/danmrichards/dessego/internal/transport"
)

// legacyMessageLimit is the limit of legacy messages to return via the API.
const legacyMessageLimit = 5

func (s *Server) getBloodMsgHandler() http.HandlerFunc {
	type getBloodMsgReq struct {
		BlockID     uint32 `form:"blockID"`
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

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		blockID := int32(bmr.BlockID)

		msgs := make([]msg.BloodMsg, 0, 10)

		// Player own messages.
		pm, err := s.ms.Player(bmr.CharacterID, blockID, bmr.ReplayNum)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining := bmr.ReplayNum - len(pm)
		msgs = append(msgs, pm...)

		// Other player messages.
		opm, err := s.ms.NonPlayer(bmr.CharacterID, blockID, remaining)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining -= len(opm)
		msgs = append(msgs, opm...)

		// Legacy messages.
		if len(msgs) < legacyMessageLimit && remaining > 0 {
			lm, err := s.ms.Legacy(blockID, bmr.ReplayNum)
			if err != nil {
				s.l.Err(err).Msg("")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			msgs = append(msgs, lm...)
		}

		// TODO: Load area names into memory for nicer logging
		s.l.Debug().Msgf(
			"%d blood messages block: %d character: %q",
			len(msgs), blockID, bmr.CharacterID,
		)

		// Message bytes.
		mb := new(bytes.Buffer)
		for _, m := range msgs {
			mb.Write(m.Bytes())
		}

		res := new(bytes.Buffer)
		binary.Write(res, binary.LittleEndian, uint32(mb.Len()))
		res.Write(mb.Bytes())

		if err = transport.WriteResponse(w, 0x1f, res.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
