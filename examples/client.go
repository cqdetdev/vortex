package main

import (
	"bytes"
	"fmt"

	"github.com/vortex-service/vortex/vortex/proto"
	"github.com/vortex-service/vortex/vortex/proto/packet"
	"github.com/gorilla/websocket"
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
	id := packet.IDLogin
	w.Varuint32(&id)
	login.Marshal(w)

	fmt.Printf("SENDING PACKET %v\n", buf.Bytes())

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
