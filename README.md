# Vortex

Websocket based microservice made in Go for simplicity and concurrency

# Mockups

```go
package main

func main() {
	s := vortex.NewService("user-database", auth.WithPassword("SECRET").WithIPWhitelist("12.34.56.78"))
	s.RegisterPackets(&FetchUserData{}, &UserData{}, &CheckAuthState{}, &AuthStateResponse{})

	db := database.New()

	s.RegisterHandle(&Handler{service: s, db: db})
	s.Start()
}

type Handler struct {
	service *AuthService
	db *Database
}

func (h *Handler) HandlePacket(conn *vortex.Conn, pk *packet.Packet) {
	switch pk {
		case *packet.FetchUserData:
			user := h.db.FetchUserData(pk.Name)
			conn.WritePacket(&UserData{
				SkinData: user.SkinData,
				RequestID: pk.RequestID
			}, false)
	}
}
```

# TODO
- Proper connection/authentication handling

# Notes

-   Uses [gophertunnel](https://github.com/sandertv/gophertunnel) packet encoding/decoding IO
-   Websockets to promote more livetime data communications
-   Possibly change to QUIC to do non-blocking async IO
