package main

import (
	"github.com/thuta/ggpoker/p2p"
)

func main() {
	// for i := 0; i < 10; i++ {
	// 	d := deck.New()
	// 	fmt.Println(d)
	// 	fmt.Println("-------")
	// }

	cfg := p2p.ServerConfig{
		ListenAddr: ":3000",
	}

	server := p2p.NewServer(cfg)
	server.Start()
}
