package packet

import (
	"github.com/cqdetdev/vortex/vortex/proto"
)

const (
	AuthResponseSuccess uint32 = iota
	AuthResponseInvalidToken
)

type AuthResponse struct {
	Code uint32
}

func (a *AuthResponse) ID() uint32 {
	return IDAuthResponse
}

func (a *AuthResponse) Marshal(io proto.IO) {
	io.Varuint32(&a.Code)
}
