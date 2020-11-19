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
			Msg("")

		// Throw a panic here as a 404 or any other form of the server accepting
		// the request will not cause Demon's Souls to treat the request as
		// failed. Need the request to be treated as failed to ensure it moves
		// on to the next one.
		panic("unhandled request")
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
		routePrefix+"/addQWCData.spd",
		middleware.LogRequest(s.l, s.addCharacterTendencyHandler()),
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
	s.r.HandleFunc(
		routePrefix+"/setWanderingGhost.spd",
		middleware.LogRequest(s.l, s.setGhostHandler()),
	)

	// Blood message routes.
	s.r.HandleFunc(
		routePrefix+"/getBloodMessage.spd",
		middleware.LogRequest(s.l, s.getBloodMsgHandler()),
	)

	// Replay routes.
	s.r.HandleFunc(
		routePrefix+"/getReplayList.spd",
		middleware.LogRequest(s.l, s.replayListHandler()),
	)
	s.r.HandleFunc(
		routePrefix+"/getReplayData.spd",
		middleware.LogRequest(s.l, s.replayDataHandler()),
	)
}
