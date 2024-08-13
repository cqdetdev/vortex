package main

import (
	"github.com/vortex-service/vortex/vortex"
	"github.com/vortex-service/vortex/vortex/auth"
	"github.com/vortex-service/vortex/vortex/proto"
	"github.com/vortex-service/vortex/vortex/proto/packet"
)

func main() {
	s := vortex.NewService("database", auth.Auth{
		Token: "TOKEN123",
	})
	s.RegisterPackets(&pingPacket{})
	s.RegisterHandler(&Handler{})
	s.Start()
}

type pingPacket struct {
	packet.Packet
}

type Handler struct{}

func (h *Handler) HandlePacket(conn *vortex.Conn, pk packet.Packet) {
	switch pk.(type) {
	case *pingPacket:
		conn.WritePacket(&pongPacket{}, false)
	}
}

func (p *pingPacket) ID() uint32 {
	return 100
}

func (p *pingPacket) Marshal(proto.IO) {}

type pongPacket struct{}

func (p *pongPacket) ID() uint32 {
	return 101
}

func (p *pongPacket) Marshal(proto.IO) {}
