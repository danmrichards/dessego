package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/danmrichards/dessego/internal/crypto"
	"github.com/danmrichards/dessego/internal/database"
	"github.com/danmrichards/dessego/internal/server/bootstrap"
	"github.com/danmrichards/dessego/internal/server/game"
	"github.com/danmrichards/dessego/internal/service/character"
	"github.com/danmrichards/dessego/internal/service/gamestate"
	"github.com/danmrichards/dessego/internal/service/ghost"
	"github.com/danmrichards/dessego/internal/service/msg"
	"github.com/danmrichards/dessego/internal/service/replay"
	"github.com/danmrichards/dessego/internal/service/sos"
	"github.com/rs/zerolog"
)

const (
	// TODO: Make configurable
	hostGame      = "127.0.0.1"
	portBootstrap = "18000"
	portUS        = "18666"
	portEU        = "18667"
	portJP        = "18668"

	dbPath = "./db/dessego.db"
)

var (
	gameServers = map[string]string{
		"US": portUS,
		"EU": portEU,
		"JP": portJP,
	}

	seed bool
)

func main() {
	flag.BoolVar(&seed, "seed", false, "Seed database tables with legacy data")
	flag.Parse()

	l := zerolog.New(os.Stdout)

	db, err := database.NewSQLite(dbPath)
	if err != nil {
		fatal(l, err)
	}
	defer db.Close()

	// Track the servers, so we can close them down later.
	servers := make([]io.Closer, 0, 4)

	// Bootstrap server; used to allow Demon's Souls to configure it's network
	// client.
	var bs *bootstrap.Server
	bs, err = bootstrap.NewServer(portBootstrap, hostGame, gameServers, l)
	if err != nil {
		fatal(l, err)
	}
	servers = append(servers, bs)

	l.Info().Msg("bootstrap server listening on " + portBootstrap)
	go func() {
		if err = bs.Serve(); err != nil {
			fatal(l, err)
		}
	}()

	// Dependencies for the gamestate server.
	rd, err := crypto.NewDecrypter(crypto.DefaultAESKey)
	if err != nil {
		fatal(l, err)
	}

	c, err := character.NewSQLiteService(db)
	if err != nil {
		fatal(l, err)
	}

	var mo []msg.Option
	if seed {
		mo = append(mo, msg.Seed())
	}
	ms, err := msg.NewSQLiteService(db, l, mo...)
	if err != nil {
		fatal(l, err)
	}

	var ro []replay.Option
	if seed {
		ro = append(ro, replay.Seed())
	}
	rs, err := replay.NewSQLiteService(db, l, ro...)
	if err != nil {
		fatal(l, err)
	}

	// Create a gamestate server for each supported region
	for region, port := range gameServers {
		gs, err := game.NewServer(
			port,
			rd,
			c,
			gamestate.NewMemory(),
			ms,
			ghost.NewMemory(l),
			rs,
			sos.NewManager(l),
			l,
		)
		if err != nil {
			fatal(l, err)
		}
		servers = append(servers, gs)

		l.Info().Msg(region + " transport server listening on " + port)
		go func() {
			if err = gs.Serve(); err != nil {
				fatal(l, err)
			}
		}()
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	l.Info().Msg("shutting down servers...")

	for _, s := range servers {
		if err = s.Close(); err != nil {
			l.Error().Err(err).Msg("close server")
		}
	}
}

func fatal(l zerolog.Logger, err error) {
	l.Fatal().Err(err).Msg("fatal error")
}
