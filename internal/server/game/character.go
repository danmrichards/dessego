package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/danmrichards/dessego/internal/service/character"
	"github.com/danmrichards/dessego/internal/transport"
)

// swagger:operation POST /cgi-bin/initializeCharacter.spd initializeCharacter
//
// Initialises a new character and persists it for future reference
//
// ---
// summary: Initialises a new character
// tags:
// - "character"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/initCharacterReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) initCharacterHandler() http.HandlerFunc {
	// swagger:model initCharacterReq
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

// swagger:operation POST /cgi-bin/getQWCData.spd getQWCData
//
// Returns the currently stored world tendency data for a given character
//
// ---
// summary: Get world tendency
// tags:
// - "character"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/worldTendencyReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) worldTendencyHandler() http.HandlerFunc {
	// swagger:model worldTendencyReq
	type worldTendencyReq struct {
		MaxNum  int `form:"maxNum"`
		Version int `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var ctr worldTendencyReq
		if err = transport.DecodeRequest(s.rd, b, &ctr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wts, err := s.cs.WorldTendency(ctr.MaxNum)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		avg := averageWorldTendency(wts)

		s.l.Debug().Msgf("current average world tendency: %q", avg)

		data := new(bytes.Buffer)
		binary.Write(data, binary.LittleEndian, avg.WB1)
		binary.Write(data, binary.LittleEndian, avg.LR1)
		binary.Write(data, binary.LittleEndian, avg.WB2)
		binary.Write(data, binary.LittleEndian, avg.LR2)
		binary.Write(data, binary.LittleEndian, avg.WB3)
		binary.Write(data, binary.LittleEndian, avg.LR3)
		binary.Write(data, binary.LittleEndian, avg.WB4)
		binary.Write(data, binary.LittleEndian, avg.LR4)
		binary.Write(data, binary.LittleEndian, avg.WB5)
		binary.Write(data, binary.LittleEndian, avg.LR5)
		binary.Write(data, binary.LittleEndian, avg.WB6)
		binary.Write(data, binary.LittleEndian, avg.LR6)
		binary.Write(data, binary.LittleEndian, avg.WB7)
		binary.Write(data, binary.LittleEndian, avg.LR7)

		if err = transport.WriteResponse(
			w, transport.ResponseCharacterTendency, data.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// swagger:operation POST /cgi-bin/addWorldTendencyHandler.spd addWorldTendencyHandler
//
// Updates the currently stored world tendency data for a given character
//
// ---
// summary: Add world tendency
// tags:
// - "character"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/addWorldTendencyReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) addWorldTendencyHandler() http.HandlerFunc {
	// swagger:model addWorldTendencyReq
	type addWorldTendencyReq struct {
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

		var atr addWorldTendencyReq
		if err = transport.DecodeRequest(s.rd, b, &atr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wt := character.WorldTendency{
			Area1: atr.Area1,
			WB1:   atr.WB1,
			LR1:   atr.LR1,
			Area2: atr.Area2,
			WB2:   atr.WB2,
			LR2:   atr.LR2,
			Area3: atr.Area3,
			WB3:   atr.WB3,
			LR3:   atr.LR3,
			Area4: atr.Area4,
			WB4:   atr.WB4,
			LR4:   atr.LR4,
			Area5: atr.Area5,
			WB5:   atr.WB5,
			LR5:   atr.LR5,
			Area6: atr.Area6,
			WB6:   atr.WB6,
			LR6:   atr.LR6,
			Area7: atr.Area7,
			WB7:   atr.WB7,
			LR7:   atr.LR7,
		}

		if err = s.cs.SetTendency(atr.CharacterID, wt); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = transport.WriteResponse(
			w, transport.ResponseAddQWCData, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// swagger:operation POST /cgi-bin/getMultiPlayGrade.spd getMultiPlayGrade
//
// Gets the current multiplayer grade for a character, used when matchmaking
//
// ---
// summary: Get multiplayer grade
// tags:
// - "character"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/multiplayerGradeReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) characterMPGradeHandler() http.HandlerFunc {
	// swagger:model multiplayerGradeReq
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

// swagger:operation POST /cgi-bin/getBloodMessageGrade.spd getBloodMessageGrade
//
// Gets the current blood message grade for a character
//
// ---
// summary: Get blood message grade
// tags:
// - "character"
// consumes:
// - text/plain
// produces:
// - text/plain
// parameters:
// - in: "body"
//   name: "body"
//   required: true
//   schema:
//     "$ref": "#/definitions/bloodMsgGradeReq"
// responses:
//   '200':
//     description: successful operation
//   '500':
//     description: unsuccessful operation
func (s *Server) characterBloodMsgGradeHandler() http.HandlerFunc {
	// swagger:model bloodMsgGradeReq
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

func averageWorldTendency(wts []character.WorldTendency) character.WorldTendency {
	at := character.WorldTendency{}
	n := len(wts)
	if n == 0 {
		return at
	}

	for _, wt := range wts {
		at.WB1 += wt.WB1
		at.WB2 += wt.WB2
		at.WB3 += wt.WB3
		at.WB4 += wt.WB4
		at.WB5 += wt.WB5
		at.WB6 += wt.WB6
		at.WB7 += wt.WB7
	}
	at.WB1 /= n
	at.WB2 /= n
	at.WB3 /= n
	at.WB4 /= n
	at.WB5 /= n
	at.WB6 /= n
	at.WB7 /= n

	return at
}
