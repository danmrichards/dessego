package game

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/service/gamestate"

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

		// Character own messages.
		cm, err := s.ms.Character(bmr.CharacterID, blockID, bmr.ReplayNum)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining := bmr.ReplayNum - len(cm)
		msgs = append(msgs, cm...)

		// Other character messages.
		ocm, err := s.ms.NonCharacter(bmr.CharacterID, blockID, remaining)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining -= len(ocm)
		msgs = append(msgs, ocm...)

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

		s.l.Debug().Msgf(
			"found %d blood messages for block: %q character: %q",
			len(msgs), gamestate.Block(blockID), bmr.CharacterID,
		)

		// Message bytes.
		mb := new(bytes.Buffer)
		for _, m := range msgs {
			mb.Write(m.Bytes())
		}

		// Response contains a header indicating the number of messages, then
		// followed by the serialised messages themselves.
		res := new(bytes.Buffer)
		binary.Write(res, binary.LittleEndian, uint32(len(msgs)))
		res.Write(mb.Bytes())

		if err = transport.WriteResponse(w, 0x1f, res.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
