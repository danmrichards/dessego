package bootstrap

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"text/template"
)

type bootstrapData struct {
	Servers map[int]string

	// TODO: Configure other bootstrap elements
}

// swagger:operation GET / bootstrap
//
// Returns a base64 encoded XML document to configure the game client for
// multiplayer.
//
// ---
// summary: Bootstrap the game client
// tags:
// - "bootstrap"
// produces:
// - text/plain
// responses:
//   '200':
//     description: successful operation
func (s *Server) bootstrapHandler() http.HandlerFunc {
	tpl := template.Must(template.ParseFiles("internal/server/bootstrap/res.tpl"))

	// Regional gamestate server URLs.
	urlUS := "http://" + s.gsHost + ":" + s.gs["US"] + "/cgi-bin/"
	urlEU := "http://" + s.gsHost + ":" + s.gs["EU"] + "/cgi-bin/"
	urlJP := "http://" + s.gsHost + ":" + s.gs["JP"] + "/cgi-bin/"

	// Data needed for gamestate bootstrap.
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
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, bd); err != nil {
			s.l.Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b := buf.Bytes()

		res := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
		base64.StdEncoding.Encode(res, b)

		w.Write(res)
	}
}
