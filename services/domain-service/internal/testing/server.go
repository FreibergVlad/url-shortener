package testing

import (
	"syscall"

	"github.com/FreibergVlad/url-shortener/domain-service/internal/config"
	"github.com/FreibergVlad/url-shortener/domain-service/internal/server"
	domainServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/domains/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration/mocks"
	"google.golang.org/grpc/test/bufconn"
)

const buffConnBufSizeBytes = 1024 * 1024

type Server struct {
	DomainServiceClient domainServiceProto.DomainServiceClient
}

func BootstrapServer() (*Server, func()) {
	config := config.New()
	listener := bufconn.Listen(buffConnBufSizeBytes)
	server := server.Bootstrap(listener, mocks.NewMockedPermissionServiceClient(true), config)

	serverStopped := make(chan struct{}, 1)
	teardown := func() {
		server.Quit <- syscall.SIGINT
		<-serverStopped
	}

	go func() {
		server.Run()
		serverStopped <- struct{}{}
	}()

	return &Server{
		DomainServiceClient: domainServiceProto.NewDomainServiceClient(integration.BufnetGrpcClient(listener)),
	}, teardown
}
