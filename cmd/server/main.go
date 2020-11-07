package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danmrichards/dessego/internal/service/gamestate"

	"github.com/danmrichards/dessego/internal/crypto"
	"github.com/danmrichards/dessego/internal/server/bootstrap"
	"github.com/danmrichards/dessego/internal/server/game"
	"github.com/danmrichards/dessego/internal/service/player"
	"github.com/danmrichards/dessego/internal/transport"
	_ "github.com/mattn/go-sqlite3"
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

var gameServers = map[string]string{
	"US": portUS,
	"EU": portEU,
	"JP": portJP,
}

func main() {
	// TODO: Logging

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if _, err = os.Create("./db/dessego.db"); err != nil {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", "./db/dessego.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Track the servers, so we can close them down later.
	servers := make([]io.Closer, 0, 4)

	// Bootstrap server; used to allow Demon's Souls to configure it's network
	// client.
	var bs *bootstrap.Server
	bs, err = bootstrap.NewServer(portBootstrap, hostGame, gameServers)
	if err != nil {
		log.Fatal(err)
	}
	servers = append(servers, bs)

	log.Println("bootstrap server listening on", portBootstrap)
	go func() {
		if err = bs.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	// Dependencies for the gamestate server.
	var rd transport.RequestDecrypter
	rd, err = crypto.NewDecrypter(crypto.DefaultAESKey)
	if err != nil {
		log.Fatal(err)
	}

	var p game.Players
	p, err = player.NewSQLiteService(db)
	if err != nil {
		log.Fatal(err)
	}

	// Create a gamestate server for each supported region
	for region, port := range gameServers {
		gs, err := game.NewServer(port, rd, p, gamestate.NewMemory())
		if err != nil {
			log.Fatal(err)
		}
		servers = append(servers, gs)

		log.Println(region+" transport server listening on", port)
		go func() {
			if err = gs.Serve(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println()
	fmt.Println("shutting down servers...")

	for _, s := range servers {
		if err = s.Close(); err != nil {
			log.Println("close server:", err)
		}
	}
}
