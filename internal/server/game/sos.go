package game

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/danmrichards/dessego/internal/service/sos"
	"github.com/danmrichards/dessego/internal/transport"
)

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
			known   = make([]uint32, 0, gsr.SOSNum)
			unknown = make([]sos.SOS, 0, gsr.SOSNum)
		)

		for _, s := range s.sos.Get(blockID, gsr.SOSNum) {
			if inSosList(strconv.FormatUint(uint64(s.ID), 10), sl) {
				known = append(known, s.ID)
			} else {
				unknown = append(unknown, s)
			}
		}

		// Response contains a header indicating the number of SOS, then
		// followed by the serialised SOS themselves.
		res := new(bytes.Buffer)

		binary.Write(res, binary.LittleEndian, uint32(len(known)))
		for _, k := range known {
			binary.Write(res, binary.LittleEndian, k)
		}

		binary.Write(res, binary.LittleEndian, uint32(len(unknown)))
		for _, u := range unknown {
			res.Write(u.Bytes())
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

func inSosList(needle string, haystack []string) bool {
	for _, h := range haystack {
		if needle == h {
			return true
		}
	}

	return false
}
