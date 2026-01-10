package p2p

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

// Server Config struct
type ServerConfig struct {
	Version     string      // App version
	ListenAddr  string      // listen address
	GameVariant GameVariant // variant of games
}

// Server struct
type Server struct {
	ServerConfig

	transport *TCPTransport
	peers     map[net.Addr]*Peer // peers map
	addPeer   chan *Peer         // add peer channel
	delPeer   chan *Peer         // delete peer channel
	msgCh     chan *Message      // message channel
}

// Initialize a new Server
func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer

	return s
}

// Start the Server
//
// telnet localhost 3000
func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("started new game server")

	if err := s.transport.ListenAndAccept(); err != nil {
		log.Fatal("listen and accept error: ", err)
	}
}

// Connect the Server
// TODO: right now we have some redundent code in registering new peers to the game network. maybe construct a new peer and handshake protocol after registering a plain connection ?
func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer

	return peer.Send([]byte(s.Version))
}

// Send the HandShake to the Peers
func (s *Server) SendHandShake(p *Peer) error {
	hs := &HandShake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}

	return p.Send(buf.Bytes())
}
