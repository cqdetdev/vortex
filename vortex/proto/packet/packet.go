package packet

import "github.com/cqdetdev/vortex/vortex/proto"

type Packet interface {
	ID() uint32
	Marshal(proto.IO)
}
