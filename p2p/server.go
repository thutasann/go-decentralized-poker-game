package p2p

import (
	"fmt"
	"net"
)

// Server Config struct
type ServerConfig struct {
	ListenAddr string // listen address
}

// Peer struct
type Peer struct {
	conn net.Conn // peer connection
}

// Server struct
type Server struct {
	ServerConfig

	listener net.Listener // net listener
	// mu       sync.RWMutex       // mutex
	peers   map[net.Addr]*Peer // peers map
	addPeer chan *Peer         // add peer channel
}

// Initialize a new Server
func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
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

	fmt.Printf("game server running on port %s", s.ListenAddr)

	// accept loop
	s.acceptLoop()
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

// Loop the server
func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("New Player connected: %s", peer.conn.RemoteAddr())
		default:
			return
		}
	}
}

// Accept the Server Loop
func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		s.handleConn(conn)
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
		fmt.Println(string(buf[:n]))
	}
}
