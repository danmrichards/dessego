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

		// Create the character, if it does not exist, in the DB.
		if err = s.cs.EnsureCreate(ucID); err != nil {
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

		s.l.Debug().Msgf("character %q logged in", ucID)

		// Response contains the character ID followed by a zero byte terminator.
		data := new(bytes.Buffer)
		data.WriteString(ucID)
		data.WriteByte(0x00)

		if err = transport.WriteResponse(
			w, transport.ResponseGeneric, data.Bytes(),
		); err != nil {
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

		ct, err := s.cs.DesiredTendency(p)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("character %q desired tendency %d", p, ct)

		// No idea why this has to be written 7 times...
		// TODO: Return real world tendency data here instead.
		data := new(bytes.Buffer)
		for i := 0; i < 7; i++ {
			for j := range []int32{int32(ct), 0} {
				binary.Write(data, binary.LittleEndian, j)
			}
		}

		if err = transport.WriteResponse(
			w, transport.ResponseCharacterTendency, data.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) addCharacterTendencyHandler() http.HandlerFunc {
	type addCharacterTendencyReq struct {
		CharacterID string `form:"characterID"`
		Area1       int    `form:"area1"`
		WB1         int    `form:"wb1"`
		LR1         int    `form:"lr1"`
		Area2       int    `form:"area2"`
		WB2         int    `form:"wb2"`
		LR2         int    `form:"lr2"`
		Area3       int    `form:"area3"`
		WB3         int    `form:"wb3"`
		LR3         int    `form:"lr3"`
		Area4       int    `form:"area4"`
		WB4         int    `form:"wb4"`
		LR4         int    `form:"lr4"`
		Area5       int    `form:"area5"`
		WB5         int    `form:"wb5"`
		LR5         int    `form:"lr5"`
		Area6       int    `form:"area6"`
		WB6         int    `form:"wb6"`
		LR6         int    `form:"lr6"`
		Area7       int    `form:"area7"`
		WB7         int    `form:"wb7"`
		LR7         int    `form:"lr7"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var atr addCharacterTendencyReq
		if err = transport.DecodeRequest(s.rd, b, &atr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Save character world tendency

		if err = transport.WriteResponse(
			w, transport.ResponseAddQWCData, []byte{0x01},
		); err != nil {
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

		stats, err := s.cs.Stats(mgr.CharacterID)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("character %q stats %s", mgr.CharacterID, stats)

		data := new(bytes.Buffer)
		for _, s := range stats.Vals() {
			binary.Write(data, binary.LittleEndian, int32(s))
		}

		if err = transport.WriteResponse(
			w, transport.ResponseCharacterMPGrade, data.Bytes(),
		); err != nil {
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

		mr, err := s.cs.MsgRating(bmr.CharacterID)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Debug().Msgf("character %q blood msg rating %d", bmr.CharacterID, mr)

		data := new(bytes.Buffer)
		binary.Write(data, binary.LittleEndian, int32(mr))

		if err = transport.WriteResponse(
			w, transport.ResponseCharacterBloodMsgGrade, data.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
