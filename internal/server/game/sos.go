package game

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danmrichards/dessego/internal/service/sos"
	"github.com/danmrichards/dessego/internal/transport"
)

type addSosDataReq struct {
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
	PlayerInfo   string  `form:"playerInfo"`
	QWCWB        uint32  `form:"qwcwb"`
	QWCLR        uint32  `form:"qwclr"`
	Black        byte    `form:"isBlack"`
	PlayerLevel  uint32  `form:"playerLevel"`
	Version      int     `form:"ver"`
}

func (a addSosDataReq) ToSos() *sos.SOS {
	return &sos.SOS{
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
		PlayerInfo:   a.PlayerInfo,
		QWCWB:        a.QWCWB,
		QWCLR:        a.QWCLR,
		Black:        a.Black,
		PlayerLevel:  a.PlayerLevel,
		Updated:      time.Now(),

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		BlockID: int32(a.BlockID),
	}
}

func (s *Server) getSosDataHandler() http.HandlerFunc {
	type getSosDataReq struct {
		BlockID        uint32 `form:"blockID"`
		MaxSOSNum      int    `form:"maxSosNum"`
		Black          int    `form:"Black"`
		Invate         int    `form:"Invate"`
		SOSNum         int    `form:"sosNum"`
		SOSList        string `form:"sosList"`
		PlayerLevelMax int    `form:"playerLevelMax"`
		PlayerLevelMin int    `form:"playerLevelMin"`
		BlackMax       int    `form:"BlackMax"`
		BlackMin       int    `form:"BlackMin"`
		InvateMax      int    `form:"InvateMax"`
		InvateMin      int    `form:"InvateMin"`
		Version        int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var gsr getSosDataReq
		if err = transport.DecodeRequest(s.rd, b, &gsr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Demon's Souls doesn't send signed integers for block IDs for some
		// reason. Coerce it.
		blockID := int32(gsr.BlockID)

		// Split the list of SOS's
		sl := strings.Split(gsr.SOSList, "a0a")

		// The client will already know about some SOS, we only need to return
		// full details for new ones.
		var (
			known   = make([]int32, 0, gsr.SOSNum)
			unknown = make([]*sos.SOS, 0, gsr.SOSNum)
		)

		for _, bs := range s.sos.List(blockID, gsr.SOSNum) {
			if inSosList(strconv.FormatUint(uint64(bs.ID), 10), sl) {
				known = append(known, bs.ID)
			} else {
				unknown = append(unknown, bs)
			}
		}

		// Response contains a header indicating the number of SOS, then
		// followed by the serialised SOS themselves.
		res := new(bytes.Buffer)

		binary.Write(res, binary.LittleEndian, uint32(len(known)))
		for _, k := range known {
			binary.Write(res, binary.LittleEndian, uint32(k))
		}

		binary.Write(res, binary.LittleEndian, uint32(len(unknown)))
		for _, u := range unknown {
			res.Write(u.Bytes())
		}

		if err = transport.WriteResponse(
			w, transport.ResponseGetSOSData, res.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) addSosDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var asr addSosDataReq
		if err = transport.DecodeRequest(s.rd, b, &asr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ns := asr.ToSos()

		// Populate the SOS with the stats for the player.
		stats, err := s.cs.Stats(asr.CharacterID)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ns.Ratings = []int{
			stats.GradeS, stats.GradeA, stats.GradeB, stats.GradeC, stats.GradeD,
		}
		ns.TotalSessions = stats.Sessions

		s.sos.Add(ns)

		if err = transport.WriteResponse(
			w, transport.ResponseAddSummonSOSData, []byte{0x01},
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) checkSosDataHandler() http.HandlerFunc {
	type checkSosDataReq struct {
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

		var csr checkSosDataReq
		if err = transport.DecodeRequest(s.rd, b, &csr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := new(bytes.Buffer)
		if rid := s.sos.Check(csr.CharacterID); rid != "" {
			data.WriteString(rid)
		} else {
			data.WriteByte(0x00)
		}

		if err = transport.WriteResponse(
			w, transport.ResponseCheckSOSData, data.Bytes(),
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) summonCharacterHandler() http.HandlerFunc {
	type summonOtherCharacterReq struct {
		GhostID  int32  `form:"ghostID"`
		NPRoomID string `form:"NPRoomID"`
		Version  int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var sor summonOtherCharacterReq
		if err = transport.DecodeRequest(s.rd, b, &sor); err != nil {
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
		s.l.Info().Msgf("player %q attempting to summon id %d", p, sor.GhostID)

		data := []byte{0x01}
		if !s.sos.Summon(sor.GhostID, sor.NPRoomID) {
			data = []byte{0x00}
			s.l.Info().Msgf(
				"player %q failed to summon non-existing id %d", p, sor.GhostID,
			)
		}

		if err = transport.WriteResponse(
			w, transport.ResponseAddSummonSOSData, data,
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) summonBlackGhostHandler() http.HandlerFunc {
	type summonBlackGhostReq struct {
		NPRoomID string `form:"NPRoomID"`
		Version  int    `form:"ver"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var sbr summonBlackGhostReq
		if err = transport.DecodeRequest(s.rd, b, &sbr); err != nil {
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
		s.l.Info().Msgf("player %q attempting to summon monk", p)

		data := []byte{0x01}
		if !s.sos.Monk(sbr.NPRoomID) {
			data = []byte{0x00}
			s.l.Info().Msgf("player %q failed to summon monk", p)
		}

		if err = transport.WriteResponse(
			w, transport.ResponseSummonMonk, data,
		); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func inSosList(needle string, haystack []string) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}

	return false
}
