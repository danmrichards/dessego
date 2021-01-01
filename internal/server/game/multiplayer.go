package game

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/danmrichards/dessego/internal/service/character"
	"github.com/danmrichards/dessego/internal/transport"
)

type finaliseMultiplayReq struct {
	CharacterID string `form:"characterID"`
	GradeS      int    `form:"gradeS"`
	GradeA      int    `form:"gradeA"`
	GradeB      int    `form:"gradeB"`
	GradeC      int    `form:"gradeC"`
	GradeD      int    `form:"gradeD"`
	Version     int    `form:"ver"`
}

func (f finaliseMultiplayReq) Grade() character.MultiplayerGrade {
	switch {
	case f.GradeS == 1:
		return character.GradeS
	case f.GradeA == 1:
		return character.GradeA
	case f.GradeB == 1:
		return character.GradeB
	case f.GradeC == 1:
		return character.GradeC
	case f.GradeD == 1:
		return character.GradeD
	default:
		return character.GradeUnknown
	}
}

type updateOtherPlayerGradeReq struct {
	CharacterID string `form:"characterID"`
	Grade       int    `form:"grade"`
	Version     int    `form:"ver"`
}

func (u updateOtherPlayerGradeReq) PlayerGrade() character.MultiplayerGrade {
	if g, ok := character.Grades[u.Grade]; ok {
		return g
	}

	return character.GradeUnknown
}

func (u updateOtherPlayerGradeReq) Character() string {
	return u.CharacterID + "0"
}

func (s *Server) outOfBlockHandler() http.HandlerFunc {
	type outOfBlockReq struct {
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

		var obr outOfBlockReq
		if err = transport.DecodeRequest(s.rd, b, &obr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.sos.Delete(obr.CharacterID)

		if err = transport.WriteResponse(
			w, transport.ResponseMultiplayerOp, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) initMultiplayHandler() http.HandlerFunc {
	type initMultiplayHandler struct {
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

		var imr initMultiplayHandler
		if err = transport.DecodeRequest(s.rd, b, &imr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = s.cs.InitMultiplayer(imr.CharacterID); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.l.Info().Msgf(
			"character %q started a multiplayer session", imr.CharacterID,
		)

		if err = transport.WriteResponse(
			w, transport.ResponseMultiplayerOp, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) finaliseMultiplayHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var fmr finaliseMultiplayReq
		if err = transport.DecodeRequest(s.rd, b, &fmr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		grade := fmr.Grade()
		if err = s.cs.UpdatePlayerGrade(fmr.CharacterID, grade); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.l.Info().Msgf(
			"character %q finished a multiplayer session and got grade %q",
			fmr.CharacterID,
			grade,
		)

		if err = transport.WriteResponse(
			w, transport.ResponseFinaliseMultiplayer, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) updateOtherPlayerGradeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var upr updateOtherPlayerGradeReq
		if err = transport.DecodeRequest(s.rd, b, &upr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		grade := upr.PlayerGrade()
		char := upr.Character()
		if err = s.cs.UpdatePlayerGrade(char, grade); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Load the current player.
		var ip string
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, err := s.gs.Player(ip)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.l.Info().Msgf(
			"character %q gave character %q grade %q", p, char, grade,
		)

		if err = transport.WriteResponse(
			w, transport.ResponseUpdateOtherPlayerGrade, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
