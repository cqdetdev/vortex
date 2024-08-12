package net

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/cqdetdev/vortex/vortex/auth"
	"github.com/cqdetdev/vortex/vortex/internal"
	"github.com/cqdetdev/vortex/vortex/proto"
	"github.com/cqdetdev/vortex/vortex/proto/packet"
	"github.com/gorilla/websocket"
)

type Server struct {
	srv *http.Server

	auth auth.Auth

	conns   []*websocket.Conn
	connsMu sync.Mutex
}

func NewServer(addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr: addr,
		},
		auth: auth.Auth{
			Token: "super-secret-token",
		},
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

		s.handle(conn)
	})

	log.Println("Server is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	return nil
}

func (s *Server) handle(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Printf("Unexpected: %v\n", err)
			break
		}

		buf := internal.BufferPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			internal.BufferPool.Put(buf)
		}()

		if len(msg) < 1 {
			return
		}

		var pk packet.Packet
		if msg[0] == byte(packet.IDLogin) {
			pk = &packet.Login{}
		}
		reader := proto.NewReader(bytes.NewReader(msg[1:]), 1, false)
		pk.Marshal(reader)

		switch pk.ID() {
		case packet.IDLogin:
			pk := pk.(*packet.Login)
			s.handleLogin(conn, pk)
		}
	}
}

func (s *Server) WritePacket(conn *websocket.Conn, pk packet.Packet, close bool) error {
	buf := internal.BufferPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		internal.BufferPool.Put(buf)
	}()

	writer := proto.NewWriter(buf, 1)
	pk.Marshal(writer)

	msg := append([]byte{byte(pk.ID())}, buf.Bytes()...)
	if close {
		if err := conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(time.Second)); err != nil {
			log.Println("Error writing close control message:", err)
			return err
		}
	} else {
		if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			log.Println("Error writing message:", err)
			return err
		}
	}

	return nil
}
