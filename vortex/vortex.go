package vortex

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/vortex-service/vortex/vortex/auth"
	"github.com/vortex-service/vortex/vortex/internal"
	"github.com/vortex-service/vortex/vortex/proto"
	"github.com/vortex-service/vortex/vortex/proto/packet"
)

type Vortex struct {
	srv *http.Server

	name string

	handler Handler
	packets []packet.Packet

	auth auth.Auth

	conns   []*Conn
	connsMu sync.Mutex
}

func NewService(name string, auth auth.Auth) *Vortex {
	return &Vortex{
		name: name,
		auth: auth,
	}
}

func (v *Vortex) Start() error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

		v.handle(conn)
	})

	log.Println("Server is listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	return nil
}

func (v *Vortex) handle(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Printf("Unexpected: %v\n", err)
			break
		}

		c := NewConn(conn)

		buf := internal.BufferPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			internal.BufferPool.Put(buf)
		}()

		if len(msg) < 1 {
			return
		}

		var pk packet.Packet
		var registeredPk bool

		if msg[0] == byte(packet.IDLogin) {
			pk = &packet.Login{}
		}

		for _, registeredPacket := range v.packets {
			if msg[0] == byte(registeredPacket.ID()) {
				pk = registeredPacket
				registeredPk = true
				break
			}
		}

		if pk == nil {
			return
		}

		reader := proto.NewReader(bytes.NewReader(msg[1:]), 1, false)
		pk.Marshal(reader)

		if registeredPk {
			v.handler.HandlePacket(c, pk)
		} else {
			switch pk.ID() {
			case packet.IDLogin:
				pk := pk.(*packet.Login)
				v.handleLogin(c, pk)
			default:
				log.Println("Received unknown packet ID:", pk.ID())
			}
		}

	}
}

func (v *Vortex) RegisterPackets(packets ...packet.Packet) {
	v.packets = append(v.packets, packets...)
}

func (v *Vortex) RegisterHandler(handler Handler) {
	v.handler = handler
}
