package main

import (
	"bytes"
	"time"

	"github.com/cqdetdev/vortex/vortex/proto"
	"github.com/cqdetdev/vortex/vortex/proto/packet"
	"github.com/gorilla/websocket"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	go func() {
		for {
			c.ReadMessage()
		}
	}()
	
	t := time.NewTicker(time.Millisecond * 500)
	defer t.Stop()

	for range t.C {
		login := &packet.Login{
			Service: "oauth",
			Token: "03ejdui39dnfksofhvw",
		}

		buf := bytes.NewBuffer(make([]byte, 0, 1024))

		w := proto.NewWriter(buf, 1)
		login.Marshal(w)

		err := c.WriteMessage(websocket.BinaryMessage, buf.Bytes())
		if err != nil {
			panic(err)
		}
	}

}