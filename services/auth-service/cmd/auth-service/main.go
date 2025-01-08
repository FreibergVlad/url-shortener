package main

import (
	"fmt"
	"net"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/server"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
)

func main() {
	config := config.NewIdentityServiceConfig()

	listener := must.Return(net.Listen("tcp", fmt.Sprintf(":%d", config.Port)))
	server, teardown := server.Bootstrap(config, listener)

	defer teardown()

	server.Run()
}
