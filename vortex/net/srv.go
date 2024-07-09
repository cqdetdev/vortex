package net

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	srv *http.Server
	
	conns []*websocket.Conn
	connsMu sync.Mutex
}

func NewServer(addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr: addr,
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
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		fmt.Println(msg)
	}

}