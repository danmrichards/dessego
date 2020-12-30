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
			lm, err := s.ms.Legacy(blockID, remaining)
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

		// Response contains a header indicating the number of messages, then
		// followed by the serialised messages themselves.
		res := new(bytes.Buffer)
		binary.Write(res, binary.LittleEndian, uint32(len(msgs)))
		for _, m := range msgs {
			res.Write(m.Bytes())
		}

		if err = transport.WriteResponse(
			w, transport.ResponseListData, res.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) addBloodMsgHandler() http.HandlerFunc {
	type addBloodMsgReq struct {
		CharacterID  string  `form:"characterID"`
		BlockID      uint32  `form:"blockID"`
		PosX         float32 `form:"posx"`
		PosY         float32 `form:"posy"`
		PosZ         float32 `form:"posz"`
		AngX         float32 `form:"angz"`
		AngY         float32 `form:"angy"`
		AngZ         float32 `form:"angz"`
		MsgID        uint32  `form:"messageID"`
		MainMsgID    uint32  `form:"mainMsgID"`
		AddMsgCateID uint32  `form:"addMsgCateID"`
		Version      int     `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var amr addBloodMsgReq
		if err = transport.DecodeRequest(s.rd, b, &amr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bm := msg.BloodMsg{
			CharacterID:  amr.CharacterID,
			BlockID:      int32(amr.BlockID),
			PosX:         amr.PosX,
			PosY:         amr.PosY,
			PosZ:         amr.PosZ,
			AngX:         amr.AngX,
			AngY:         amr.AngY,
			AngZ:         amr.AngZ,
			MsgID:        amr.MsgID,
			MainMsgID:    amr.MainMsgID,
			AddMsgCateID: amr.AddMsgCateID,
		}

		if err = s.ms.Add(bm); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("added new message %q", bm)

		if err = transport.WriteResponse(
			w, transport.ResponseAddData, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) deleteBloodMsgHandler() http.HandlerFunc {
	type deleteBloodMsgReq struct {
		BloodMsgID int `form:"bmID"`
		Version    int `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var dmr deleteBloodMsgReq
		if err = transport.DecodeRequest(s.rd, b, &dmr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = s.ms.Delete(dmr.BloodMsgID); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("deleted message %d", dmr.BloodMsgID)

		if err = transport.WriteResponse(
			w, transport.ResponseDeleteBloodMsg, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) updateBloodMsgGradeHandler() http.HandlerFunc {
	type updateBloodMsgGradeReq struct {
		BloodMsgID int `form:"bmID"`
		Version    int `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var ugr updateBloodMsgGradeReq
		if err = transport.DecodeRequest(s.rd, b, &ugr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		bm, err := s.ms.Get(ugr.BloodMsgID)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = s.ms.UpdateRating(ugr.BloodMsgID); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.l.Debug().Msgf("recommended message %q", bm)

		if err = s.cs.UpdateMsgRating(bm.CharacterID); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.l.Debug().Msgf(
			"updated message rating for character: %q", bm.CharacterID,
		)

		if err = transport.WriteResponse(
			w, transport.ResponseUpdateMsgGrade, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
