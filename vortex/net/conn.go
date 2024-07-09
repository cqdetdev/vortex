package net

import (
	"github.com/cqdetdev/vortex/vortex/proto"
	"github.com/gorilla/websocket"
)

type Conn struct {
	conn *websocket.Conn

	dec *proto.Decoder
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{
		conn: conn,
		dec:  proto.NewDecoder(conn),
	}
}