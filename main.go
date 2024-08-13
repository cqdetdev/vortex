package main

import (
	"github.com/vortex-service/vortex/vortex"
	"github.com/vortex-service/vortex/vortex/auth"
)

func main() {
	s := vortex.NewService("database", auth.Auth{
		Token: "TOKEN123",
	})
	s.Start()
}
