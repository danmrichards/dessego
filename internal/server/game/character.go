package game

import (
	"fmt"
	"io/ioutil"
	"log"
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
		log.Printf(
			"init character request from %q to %q", r.RemoteAddr, s.l.Addr(),
		)

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var icr initCharacterReq
		if err = transport.DecodeRequest(s.rd, b, &icr); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Printf("%+v\n", icr)

		// TODO: Create player in DB
		// TODO: Add player to active list
	}
}
