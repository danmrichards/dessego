package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danmrichards/dessego/internal/server/bootstrap"
	"github.com/danmrichards/dessego/internal/server/game"
)

const (
	// TODO: Make configurable
	host          = "127.0.0.1"
	portBootstrap = "18000"
	portUS        = "18666"
	portEU        = "18667"
	portJP        = "18668"
)

var gameServers = map[string]string{
	"US": portUS,
	"EU": portEU,
	"JP": portJP,
}

func main() {
	// TODO: Logging

	servers := make([]io.Closer, 0, 4)

	bs, err := bootstrap.NewServer(host, portBootstrap, gameServers)
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

	for region, port := range gameServers {
		gs, err := game.NewServer(host, port)
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
