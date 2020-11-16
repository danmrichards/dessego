package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
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

		// Unique character ID.
		ucID := fmt.Sprintf("%s%d", icr.CharacterID, icr.Index)

		// Create the player, if it does not exist, in the DB.
		if err = s.ps.EnsureCreate(ucID); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var ip string
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Track the player in game state.
		s.gs.AddPlayer(ip, ucID)

		data := new(bytes.Buffer)
		data.WriteString(ucID)
		data.WriteByte(0x00)

		if err = transport.WriteResponse(w, 0x17, data.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) characterTendencyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p, err := s.gs.Player(ip)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ct, err := s.ps.DesiredTendency(p)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("player %q desired tendency %d", p, ct)

		// No idea why this has to be written 7 times...
		data := new(bytes.Buffer)
		for i := 0; i < 7; i++ {
			for j := range []int32{int32(ct), 0} {
				binary.Write(data, binary.LittleEndian, j)
			}
		}

		if err = transport.WriteResponse(w, 0x0e, data.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) characterMPGradeHandler() http.HandlerFunc {
	type multiplayerGradeReq struct {
		CharacterID string `form:"NPID"`
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

		var mgr multiplayerGradeReq
		if err = transport.DecodeRequest(s.rd, b, &mgr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stats, err := s.ps.Stats(mgr.CharacterID)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("player %q stats %+v", mgr.CharacterID, stats)

		data := new(bytes.Buffer)
		for _, s := range stats.Vals() {
			binary.Write(data, binary.LittleEndian, int32(s))
		}

		if err = transport.WriteResponse(w, 0x28, data.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) characterBloodMsgGradeHandler() http.HandlerFunc {
	type bloodMsgGradeReq struct {
		CharacterID string `form:"NPID"`
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

		var bmr bloodMsgGradeReq
		if err = transport.DecodeRequest(s.rd, b, &bmr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mr, err := s.ps.MsgRating(bmr.CharacterID)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("player %q blood msg rating %d", bmr.CharacterID, mr)

		data := new(bytes.Buffer)
		binary.Write(data, binary.LittleEndian, int32(mr))

		if err = transport.WriteResponse(w, 0x29, data.Bytes()); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
