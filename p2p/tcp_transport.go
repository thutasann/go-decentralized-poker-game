package p2p

import (
	"net"

	"github.com/sirupsen/logrus"
)

// TCP Transport struct
type TCPTransport struct {
	listenAddr string // listen address
	listener   net.Listener
	AddPeer    chan *Peer
	DelPeer    chan *Peer
}

// Initialize a new TCP Transport
func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}

// Listen and Accept
func (t *TCPTransport) ListenAndAccept() error {
	listen, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = listen

	for {
		conn, err := listen.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}

		peer := &Peer{
			conn: conn,
		}

		t.AddPeer <- peer
	}

}
