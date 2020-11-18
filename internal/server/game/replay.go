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

func (s *Server) replayListHandler() http.HandlerFunc {
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
			"found %d replays for block: %q", len(lr), gamestate.Block(blockID),
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

func (s *Server) replayDataHandler() http.HandlerFunc {
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
