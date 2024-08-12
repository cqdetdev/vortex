package vortex

import "github.com/vortex-service/vortex/vortex/net"

type Vortex struct {
	srv *net.Server
}

func New() *Vortex {
	return &Vortex{
		srv: net.NewServer(":8080"),
	}
}

func (v *Vortex) Start() error {
	return v.srv.Start()
}