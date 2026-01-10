package main

import (
	"log"
	"time"

	"github.com/thuta/ggpoker/p2p"
)

func main() {
	// server config
	cfg := p2p.ServerConfig{
		Version:     "GGPOKER V0.1-beta",
		ListenAddr:  ":3000",
		GameVariant: p2p.TexasHoldem,
	}
	server := p2p.NewServer(cfg)
	go server.Start()

	// Simulate
	time.Sleep(1 * time.Second)

	// remote server config
	remoteCfg := p2p.ServerConfig{
		Version:     "GGPOKER V0.1-beta",
		ListenAddr:  ":4000",
		GameVariant: p2p.TexasHoldem,
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	if err := remoteServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}

	select {}
}
