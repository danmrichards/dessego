package game

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *Server) getBloodMsgHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		req := s.rd.Decrypt(b)
		fmt.Println(req)

		// TODO: Get Blood Message
	}
}
