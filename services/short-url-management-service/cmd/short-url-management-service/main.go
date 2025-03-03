package main

import (
	"fmt"
	"net"

	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.NewShortURLManagerConfig()
	grpcDialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())

	permissionServiceConn := must.Return(grpc.NewClient(config.AuthServiceDSN, grpcDialOptions))
	permissionServiceClient := permissionServiceProto.NewPermissionServiceClient(permissionServiceConn)

	listener := must.Return(net.Listen("tcp", fmt.Sprintf(":%d", config.Port)))
	server, teardown := server.BootstrapShortURLManagementService(config, listener, permissionServiceClient)

	defer teardown()

	server.Run()
}
