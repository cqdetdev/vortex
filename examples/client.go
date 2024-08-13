package main

import (
	"bytes"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/vortex-service/vortex/vortex/proto"
	"github.com/vortex-service/vortex/vortex/proto/packet"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	login := &packet.Login{
		Service: "oauth-service",
		Token:   "super-secret-token",
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	w := proto.NewWriter(buf, 1)
	id := uint32(100)
	w.Varuint32(&id)
	login.Marshal(w)

	fmt.Printf("SENDING LOGIN PACKET %v\n", buf.Bytes())

	err = c.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	if err != nil {
		panic(err)
	}

	ping := &pingPacket{}
	buf.Reset()
	w = proto.NewWriter(buf, 1)
	id = ping.ID()
	w.Varuint32(&id)
	ping.Marshal(w)

	fmt.Printf("SENDING PING PACKET %v\n", buf.Bytes())

	err = c.WriteMessage(websocket.BinaryMessage, buf.Bytes())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
			}

			fmt.Println(msg)
		}
	}()

	select {}
}

type pingPacket struct{}

func (p *pingPacket) ID() uint32 {
	return 100
}

func (p *pingPacket) Marshal(proto.IO) {}
