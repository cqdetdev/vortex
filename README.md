# Vortex

Websocket based microservice for whatever you wanna do

# Mockups

```go
type AuthService struct {
    // DATA
}

func(s AuthService) Start() {

}

// On data receive
func(s AuthService) Recv(pk packet.Packet) {

}
```

# Other things

-   Uses gophertunnel encoding style for packets
-   Websockets because idk they're cool
