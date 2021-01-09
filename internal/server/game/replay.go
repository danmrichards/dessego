package game

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/service/gamestate"
	"github.com/danmrichards/dessego/internal/service/replay"
	"github.com/danmrichards/dessego/internal/transport"
)

// swagger:model addReplayDataReq
type addReplayDataReq struct {
	CharacterID  string  `form:"characterID"`
	BlockID      uint32  `form:"blockID"`
	PosX         float32 `form:"posx"`
	PosY         float32 `form:"posy"`
	PosZ         float32 `form:"posz"`
	AngX         float32 `form:"angx"`
	AngY         float32 `form:"angy"`
	AngZ         float32 `form:"angz"`
	MsgID        uint32  `form:"messageID"`
	MainMsgID    uint32  `form:"mainMsgID"`
	AddMsgCateID uint32  `form:"addMsgCateID"`
	Data         string  `form:"replayBinary"`
}

func (a addReplayDataReq) ToReplay() *replay.Replay {
	return &replay.Replay{
		CharacterID:  a.CharacterID,
		PosX:         a.PosX,
		PosY:         a.PosY,
		PosZ:         a.PosZ,
		AngX:         a.AngX,
		AngY:         a.AngY,
		AngZ:         a.AngZ,
		MsgID:        a.MsgID,
		MainMsgID:    a.MainMsgID,
		AddMsgCateID: a.AddMsgCateID,
		Data:         []byte(a.Data),

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		BlockID: int32(a.BlockID),
	}
}

// swagger:operation POST /cgi-bin/getReplayList.spd replayListHandler
//
// Returns a list of available replays for an area of the game
//
// ---
// summary: List replays
// tags:
// - "replays"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/replayListReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) replayListHandler() http.HandlerFunc {
	// swagger:model replayListReq
	type replayListReq struct {
		BlockID   uint32 `form:"blockID"`
		ReplayNum int    `form:"replayNum"`
		Version   int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var rlr replayListReq
		if err = transport.DecodeRequest(s.rd, b, &rlr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		blockID := int32(rlr.BlockID)

		rs := make([]replay.Replay, 0, 10)

		// Non-legacy replays.
		nlr, err := s.rs.List(blockID, rlr.ReplayNum, replay.NonLegacy)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		remaining := rlr.ReplayNum - len(nlr)
		rs = append(rs, nlr...)

		// Legacy replays.
		lr, err := s.rs.List(blockID, remaining, replay.Legacy)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rs = append(rs, lr...)

		s.l.Debug().Msgf(
			"found %d replays for block: %q", len(rs), gamestate.Block(blockID),
		)

		// Response contains a header indicating the number of replays, then
		// followed by the serialised replay headers themselves.
		res := new(bytes.Buffer)
		binary.Write(res, binary.LittleEndian, uint32(len(rs)))
		for _, r := range rs {
			res.Write(r.Header())
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

// swagger:operation POST /cgi-bin/getReplayData.spd getReplayDataHandler
//
// Returns a single replay's data
//
// ---
// summary: Get replay data
// tags:
// - "replays"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/replayDataReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) getReplayDataHandler() http.HandlerFunc {
	// swagger:model replayDataReq
	type replayDataReq struct {
		GhostID uint32 `form:"ghostID"`
		Version int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var rdr replayDataReq
		if err = transport.DecodeRequest(s.rd, b, &rdr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rp, err := s.rs.Get(rdr.GhostID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.l.Warn().Msgf("no replay exists with ID: %d", rdr.GhostID)
		case err != nil:
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("loading replay %s", rp)

		// Response is in the format ghost ID, replay length followed by replay
		// data.
		res := new(bytes.Buffer)
		binary.Write(res, binary.LittleEndian, rdr.GhostID)
		binary.Write(res, binary.LittleEndian, uint32(len(rp.Data)))
		res.Write(rp.Data)

		if err = transport.WriteResponse(
			w, transport.ResponseReplayData, res.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// swagger:operation POST /cgi-bin/addReplayData.spd addReplayDataHandler
//
// Adds replay data for the given character
//
// ---
// summary: Add replay data
// tags:
// - "replays"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/addReplayDataReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) addReplayDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var adr addReplayDataReq
		if err = transport.DecodeRequest(s.rd, b, &adr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		nr := adr.ToReplay()
		if err = s.rs.Add(nr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("added new replay replay %s", nr)

		if err = transport.WriteResponse(
			w, transport.ResponseAddData, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
