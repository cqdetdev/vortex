package vortex

import (
	"fmt"

	"github.com/vortex-service/vortex/vortex/proto/packet"
)

type Handler interface {
	HandlePacket(conn *Conn, pk *packet.Packet)
}

func (s *Vortex) handleLogin(c *Conn, pk *packet.Login) {
	resp := &packet.AuthResponse{}
	var closed bool
	if s.auth.Token == pk.Token {
		resp.Code = packet.AuthResponseSuccess
		closed = false
		s.connsMu.Lock()
		defer s.connsMu.Unlock()
	} else {
		resp.Code = packet.AuthResponseInvalidToken
		closed = true
	}

	err := c.WritePacket(pk, closed)
	if err != nil {
		fmt.Println(err)
	}
}
