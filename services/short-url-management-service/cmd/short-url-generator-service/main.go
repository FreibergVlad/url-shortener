package main

import (
	"fmt"
	"net"

	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	userServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/users/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config := config.NewShortURLGeneratorConfig()
	grpcDialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())

	grpcPermissionServiceConn := must.Return(grpc.NewClient(config.AuthServiceDSN, grpcDialOptions))
	grpcPermissionServiceClient := permissionServiceProto.NewPermissionServiceClient(grpcPermissionServiceConn)

	grpcUserServiceConn := must.Return(grpc.NewClient(config.AuthServiceDSN, grpcDialOptions))
	grpcUserServiceClient := userServiceProto.NewUserServiceClient(grpcUserServiceConn)

	grpcDomainServiceConn := must.Return(grpc.NewClient(config.DomainServiceDSN, grpcDialOptions))
	grpcDomainServiceClient := domainServiceProto.NewDomainServiceClient(grpcDomainServiceConn)

	listener := must.Return(net.Listen("tcp", fmt.Sprintf(":%d", config.Port)))
	server, teardown := server.BootstrapShortURLGeneratorService(
		config, listener, grpcPermissionServiceClient, grpcUserServiceClient, grpcDomainServiceClient,
	)

	defer teardown()

	server.Run()
}
