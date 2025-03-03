package testing

import (
	"syscall"

	shortUrlGeneratorServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/generator/service/v1"
	shortUrlManagementServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/shorturls/management/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration/mocks"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/migrations"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/server"
	"google.golang.org/grpc/test/bufconn"
)

const (
	ShortURLDomain       = "fake-domain.com"
	buffConnBufSizeBytes = 1024 * 1024
)

type Server struct {
	ShortURLManagementServiceClient shortUrlManagementServiceProto.ShortURLManagementServiceClient
	ShortURLGeneratorServiceClient  shortUrlGeneratorServiceProto.ShortURLGeneratorServiceClient
}

func BootstrapServer() (*Server, func()) {
	migrations.Migrate(config.NewMigrationConfig())

	permissionServiceClient := mocks.NewMockedPermissionServiceClient(true)

	domainServiceClient := NewMockedDomainServiceClient([]string{ShortURLDomain})
	userServiceClient := NewMockedUserServiceClient(FakeUser)

	shortURLManagementListener := bufconn.Listen(buffConnBufSizeBytes)
	shortURLManagementServer, shortURLManagementTeardown := server.BootstrapShortURLManagementService(
		config.NewShortURLManagerConfig(),
		shortURLManagementListener,
		permissionServiceClient,
	)
	shortURLGeneratorListener := bufconn.Listen(buffConnBufSizeBytes)
	shortURLGeneratorServer, shortURLGeneratorTeardown := server.BootstrapShortURLGeneratorService(
		config.NewShortURLGeneratorConfig(),
		shortURLGeneratorListener,
		permissionServiceClient,
		userServiceClient,
		domainServiceClient,
	)

	shortURLManagerStopped, shortURLGeneratorStopped := make(chan struct{}, 1), make(chan struct{}, 1)
	teardown := func() {
		shortURLManagementServer.Quit <- syscall.SIGINT
		shortURLManagementTeardown()
		<-shortURLManagerStopped

		shortURLGeneratorServer.Quit <- syscall.SIGINT
		shortURLGeneratorTeardown()
		<-shortURLGeneratorStopped
	}

	go func() {
		shortURLManagementServer.Run()
		shortURLManagerStopped <- struct{}{}
	}()

	go func() {
		shortURLGeneratorServer.Run()
		shortURLGeneratorStopped <- struct{}{}
	}()

	return &Server{
		ShortURLManagementServiceClient: shortUrlManagementServiceProto.NewShortURLManagementServiceClient(
			integration.BufnetGrpcClient(shortURLManagementListener),
		),
		ShortURLGeneratorServiceClient: shortUrlGeneratorServiceProto.NewShortURLGeneratorServiceClient(
			integration.BufnetGrpcClient(shortURLGeneratorListener),
		),
	}, teardown
}
