package packet

import (
	"github.com/vortex-service/vortex/vortex/proto"
)

type Heartbeat struct {
	Timestamp int64
}

func (h *Heartbeat) ID() uint32 {
	return IDHeartbeat
}

func (h *Heartbeat) Marshal(io proto.IO) {
	io.Int64(&h.Timestamp)
}
