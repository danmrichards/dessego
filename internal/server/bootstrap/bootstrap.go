package bootstrap

import (
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"text/template"
)

type bootstrapData struct {
	Servers map[int]string

	// TODO: Configure other bootstrap elements
}

func (s *Server) handleBootstrap() http.HandlerFunc {
	tpl := template.Must(template.ParseFiles("internal/server/bootstrap/res.tpl"))

	// Regional game server URLs.
	urlUS := "http://" + s.host + ":" + s.gs["US"] + "/cgi-bin/"
	urlEU := "http://" + s.host + ":" + s.gs["EU"] + "/cgi-bin/"
	urlJP := "http://" + s.host + ":" + s.gs["JP"] + "/cgi-bin/"

	// Data needed for game bootstrap.
	bd := bootstrapData{
		Servers: map[int]string{
			1:  urlUS,
			2:  urlEU,
			3:  urlJP,
			4:  urlJP,
			5:  urlEU,
			6:  urlEU,
			7:  urlEU,
			8:  urlEU,
			11: urlJP,
			12: urlJP,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("bootstrap request from %q to %q", r.RemoteAddr, s.l.Addr())

		var buf bytes.Buffer
		if err := tpl.Execute(&buf, bd); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b := buf.Bytes()

		res := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
		base64.StdEncoding.Encode(res, b)

		w.Write(res)
	}
}
