package game

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/danmrichards/dessego/internal/service/ghost"

	"github.com/danmrichards/dessego/internal/service/gamestate"
	"github.com/danmrichards/dessego/internal/transport"
	dsbase64 "github.com/danmrichards/dessego/internal/transport/encoding/base64"
)

// TODO: Configurable?
const maxGhostAge = 30 * time.Second

func (s *Server) getGhostHandler() http.HandlerFunc {
	type getGhostReq struct {
		Version     int      `form:"ver"`
		CharacterID string   `form:"characterID"`
		BlockID     uint32   `form:"blockID"`
		MaxGhosts   int      `form:"maxGhostNum"`
		SOS         int      `form:"sosNum"`
		SOSIDList   []string `form:"sosIDList"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var ggr getGhostReq
		if err = transport.DecodeRequest(s.rd, b, &ggr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		blockID := int32(ggr.BlockID)

		// Clear out stale ghosts to ensure up to date replays.
		s.gh.ClearBefore(time.Now().Add(-maxGhostAge))

		g := s.gh.Get(ggr.CharacterID, blockID, ggr.MaxGhosts)
		s.l.Debug().Msgf(
			"found %d ghosts for block: %q character: %q",
			len(g), gamestate.Block(blockID), ggr.CharacterID,
		)

		// Response contains a header indicating the number of ghosts, followed
		// by the ghost replay data itself.
		res := new(bytes.Buffer)
		binary.Write(res, binary.LittleEndian, uint32(0))
		binary.Write(res, binary.LittleEndian, uint32(len(g)))

		// TODO: Encode replay data.

		if err = transport.WriteResponse(w, 0x11, res.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) setGhostHandler() http.HandlerFunc {
	type setGhostReq struct {
		CharacterID  string  `form:"characterID"`
		GhostBlockID uint32  `form:"ghostBlockID"`
		PosX         float32 `form:"posx"`
		PosY         float32 `form:"posy"`
		PosZ         float32 `form:"posz"`
		ReplayData   string  `form:"replayData"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var sgr setGhostReq
		if err = transport.DecodeRequest(s.rd, b, &sgr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		blockID := int32(sgr.GhostBlockID)

		// Cannot use the std library decoding as Demon's Souls sends replay
		// data with broken encoding.
		rd, err := dsbase64.StdEncoding.DecodeString(sgr.ReplayData)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		g := ghost.NewGhost(blockID, sgr.CharacterID, rd)

		// Check if the character has spawned or changed area.
		prev, err := s.gh.Character(sgr.CharacterID)
		if err != nil {
			var cgerr ghost.CharacterGhostNotFoundError
			if !errors.As(err, &cgerr) {
				s.l.Err(err).Msg("")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = nil

			s.l.Debug().Msgf(
				"character: %q spawned into block: %q",
				sgr.CharacterID, gamestate.Block(blockID),
			)
		} else if prev.BlockID != blockID {
			s.l.Debug().Msgf(
				"character: %q moved from block: %q to block: %q",
				sgr.CharacterID, gamestate.Block(prev.BlockID), gamestate.Block(blockID),
			)
		}

		s.gh.Set(sgr.CharacterID, g)

		if err = transport.WriteResponse(w, 0x17, []byte{0x01}); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
