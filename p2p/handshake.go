package p2p

import (
	"encoding/gob"

	"github.com/sirupsen/logrus"
)

// HandShake represents the handshake between peers
type HandShake struct {
	Version     string      // Version
	GameVariant GameVariant // Game Variant
}

// HandShake the peers
func (s *Server) handshake(p *Peer) error {
	hs := &HandShake{}
	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"peer":    p.conn.RemoteAddr(),
		"version": hs.Version,
		"variant": hs.GameVariant,
	}).Info("[handshake] received handshake")

	return nil
}
