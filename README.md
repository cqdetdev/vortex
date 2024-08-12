# Vortex

Websocket based microservice made in Go for simplicity and concurrency

# Mockups

```go
type AuthService struct {
	Token string
}

func(s AuthService) Start() {
	vortex.RegisterPackets(&FetchUserData{}, &CheckAuthState{}, &AuthStateResponse{})
	// Or ditch and use a global response packet + json + packet request/response ids
}

func(s AuthService) Recv(pk packet.Packet) {

}
```

# Notes

-   Uses [gophertunnel](https://github.com/sandertv/gophertunnel) packet encoding/decoding IO
-   Websockets to promote more livetime data communications
-   Possibly change to QUIC to do non-blocking async IO
