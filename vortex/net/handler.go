package net

import (
	"fmt"

	"github.com/cqdetdev/vortex/vortex/proto/packet"
	"github.com/gorilla/websocket"
)

func (s *Server) handleLogin(c *websocket.Conn, pk *packet.Login) {
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

	err := s.WritePacket(c, pk, closed)
	if err != nil {
		fmt.Println(err)
	}
}
