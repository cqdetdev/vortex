package vortex

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vortex-service/vortex/vortex/internal"
	"github.com/vortex-service/vortex/vortex/proto"
	"github.com/vortex-service/vortex/vortex/proto/packet"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Conn struct {
	conn *websocket.Conn
}

func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{conn: conn}
}

func (c *Conn) WritePacket(pk packet.Packet, close bool) error {
	buf := internal.BufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		internal.BufferPool.Put(buf)
	}()

	writer := proto.NewWriter(buf, 1)
	pk.Marshal(writer)

	msg := append([]byte{byte(pk.ID())}, buf.Bytes()...)
	if close {
		if err := c.conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second)); err != nil {
			log.Println("Error writing close control message:", err)
			return err
		}
	} else {
		if err := c.conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			log.Println("Error writing message:", err)
			return err
		}
	}

	return nil
}
