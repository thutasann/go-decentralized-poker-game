package p2p

import (
	"bytes"
	"io"
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

// Message Struct
type Message struct {
	Payload io.Reader
	From    net.Addr
}

// Peer struct
type Peer struct {
	conn net.Conn // peer connection
}

// Send Fn
func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// Read Loop Fn
func (p *Peer) ReadLoop(msgch chan *Message) {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		msgch <- &Message{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}
	}

	// TODO: unregister this peer
	if err := p.conn.Close(); err != nil {
		log.Fatal("errror at connection close: ", err)
	}
}

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
