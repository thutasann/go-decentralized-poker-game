package p2p

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

// Server Config struct
type ServerConfig struct {
	ListenAddr string // listen address
}

// Server struct
type Server struct {
	ServerConfig

	handler  Handler
	listener net.Listener       // net listener
	peers    map[net.Addr]*Peer // peers map
	addPeer  chan *Peer         // add peer channel
	msgCh    chan *Message      // message channel
	// mu       sync.RWMutex       // mutex
}

// Initialize a new Server
func NewServer(cfg ServerConfig) *Server {
	return &Server{
		handler:      &DefaultHandler{},
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}
}

// Start the Server
//
// telnet localhost 3000
func (s *Server) Start() {
	// loop
	go s.loop()

	// listen
	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server running on port %s\n", s.ListenAddr)

	// accept loop
	s.acceptLoop()
}

// Loop the server
func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("New Player connected: %s\n", peer.conn.RemoteAddr())
		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

// Listen the server
func (s *Server) listen() error {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		panic(err)
	}

	s.listener = listen
	return nil
}

// Accept the Server Loop
func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		s.addPeer <- peer

		if err := peer.Send([]byte("GGPOKER V0.1-beta")); err != nil {
			log.Printf("failed to send handshake to peer: %v", err)
		}

		go s.handleConn(conn)
	}
}

// Handle the Server Connection
func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		s.msgCh <- &Message{
			From:    conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}
	}
}
