# Vortex

Websocket based microservice made in Go for simplicity and concurrency

# Mockups

```go
type AuthService struct {
	Token string
	ExtraData any

	Service vortex.Service
}

func(s AuthService) Start() {
	s.Service = vortex.NewService("user-database", auth.WithPassword("SECRET").WithIPWhitelist("12.34.56.78"))
	s.Service.RegisterPackets(&FetchUserData{}, &CheckAuthState{}, &AuthStateResponse{})
	s.Service.Start()
}

type Handler struct {
	service *AuthService
}

func (h *Handler) HandlePacket(pk *packet.Packet) bool {

}
```

# Notes

-   Uses [gophertunnel](https://github.com/sandertv/gophertunnel) packet encoding/decoding IO
-   Websockets to promote more livetime data communications
-   Possibly change to QUIC to do non-blocking async IO
