package game

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/danmrichards/dessego/internal/transport"
)

func (s *Server) getBloodMsgHandler() http.HandlerFunc {
	type getBloodMsgReq struct {
		BlockID     int    `form:"blockID"`
		ReplyNum    int    `form:"replayNum"`
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

		var bmr getBloodMsgReq
		if err = transport.DecodeRequest(s.rd, b, &bmr); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v\n", bmr)

		// TODO: Get Blood Message
	}
}
