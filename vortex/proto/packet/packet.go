package packet

import "github.com/vortex-service/vortex/vortex/proto"

type Packet interface {
	ID() uint32
	Marshal(proto.IO)
}
