package p2p

import "fmt"

// Handle the Message
func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%v\n", msg)
	return nil
}
