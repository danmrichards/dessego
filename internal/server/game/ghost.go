package game

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/transport"
)

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

		fmt.Printf("block: %d req: %+v\n", blockID, ggr)

		// TODO: Get Ghost
	}
}
