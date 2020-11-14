package game

import (
	"net/http"

	"github.com/danmrichards/dessego/internal/server/middleware"
)

const routePrefix = "/cgi-bin"

func (s *Server) routes() {
	s.r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.l.Warn().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("client", r.RemoteAddr).
			Msg("unhandled request")
	})

	// System routes.
	s.r.HandleFunc(
		routePrefix+"/login.spd",
		middleware.LogRequest(s.l, s.loginHandler()),
	)
	s.r.HandleFunc(
		routePrefix+"/getTimeMessage.spd",
		middleware.LogRequest(s.l, s.timeMsgHandler()),
	)

	// Character/Player routes.
	s.r.HandleFunc(
		routePrefix+"/initializeCharacter.spd",
		middleware.LogRequest(s.l, s.initCharacterHandler()),
	)
	s.r.HandleFunc(
		routePrefix+"/getQWCData.spd",
		middleware.LogRequest(s.l, s.characterTendencyHandler()),
	)
	s.r.HandleFunc(
		routePrefix+"/getMultiPlayGrade.spd",
		middleware.LogRequest(s.l, s.characterMPGradeHandler()),
	)
	s.r.HandleFunc(
		routePrefix+"/getBloodMessageGrade.spd",
		middleware.LogRequest(s.l, s.characterBloodMsgGradeHandler()),
	)

	// Ghost routes.
	s.r.HandleFunc(
		routePrefix+"/getWanderingGhost.spd",
		middleware.LogRequest(s.l, s.getGhostHandler()),
	)

	// Blood message routes.
	s.r.HandleFunc(
		routePrefix+"/getBloodMessage.spd",
		middleware.LogRequest(s.l, s.getBloodMsgHandler()),
	)
}
