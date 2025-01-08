package main

import (
	"fmt"
	"net"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	"github.com/FreibergVlad/url-shortener/domain-service/internal/server"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.New()

	permissionServiceConn := must.Return(grpc.NewClient(config.AuthServiceDSN, grpc.WithTransportCredentials(insecure.NewCredentials())))
	permissionServiceClient := permissionServiceProto.NewPermissionServiceClient(permissionServiceConn)

	listener := must.Return(net.Listen("tcp", fmt.Sprintf(":%d", config.Port)))
	server := server.Bootstrap(listener, permissionServiceClient, config)
	server.Run()
}
