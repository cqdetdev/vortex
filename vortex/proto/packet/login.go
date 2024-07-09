package packet

import "github.com/cqdetdev/vortex/vortex/proto"

type Login struct {
	Service string
	Token   string
}

func (l *Login) ID() uint32 {
	return IDLogin
}

func (l *Login) Marshal(io proto.IO) {
	io.String(&l.Service)
	io.String(&l.Token)
}