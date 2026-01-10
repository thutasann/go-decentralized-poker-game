package p2p

import "github.com/sirupsen/logrus"

// loop runs the main server event loop.
// It serializes all server-side state mutations by handling peer lifecycle
// events (connect/disconnect) and inbound messages through channels.
//
// This design avoids explicit locking by ensuring that all access to shared
// state (e.g. s.peers) occurs on a single goroutine, making the server
// concurrency-safe by construction.
func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		case peer := <-s.addPeer:
			// Handshake
			if err := s.handshake(peer); err != nil {
				logrus.Info("handshake with incoming player failed: ", err)
				continue
			}

			// TODO: check max players and other game state logic
			go peer.ReadLoop(s.msgCh)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
