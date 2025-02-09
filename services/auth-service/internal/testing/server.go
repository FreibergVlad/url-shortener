package testing

import (
	"syscall"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/migrations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/server"
	invitationServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/invitations/service/v1"
	organizationServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/organizations/service/v1"
	permissionServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/permissions/service/v1"
	tokenServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/tokens/service/v1"
	userServiceProto "github.com/FreibergVlad/url-shortener/proto/pkg/users/service/v1"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/testing/integration"
	"google.golang.org/grpc/test/bufconn"
)

const buffConnBufSizeBytes = 1024 * 1024

type Server struct {
	UserServiceClient         userServiceProto.UserServiceClient
	TokenServiceClient        tokenServiceProto.TokenServiceClient
	PermissionServiceClient   permissionServiceProto.PermissionServiceClient
	OrganizationServiceClient organizationServiceProto.OrganizationServiceClient
	InvitationServiceClient   invitationServiceProto.InvitationServiceClient
}

func BootstrapServer() (*Server, func()) {
	migrations.Migrate(config.NewMigrationConfig())

	listener := bufconn.Listen(buffConnBufSizeBytes)
	server, baseTeardown := server.Bootstrap(config.NewIdentityServiceConfig(), listener)
	stopped := make(chan struct{}, 1)

	teardown := func() {
		server.Quit <- syscall.SIGINT
		baseTeardown()
		<-stopped
	}

	go func() {
		server.Run()
		stopped <- struct{}{}
	}()

	client := integration.BufnetGrpcClient(listener)
	return &Server{
		UserServiceClient:         userServiceProto.NewUserServiceClient(client),
		TokenServiceClient:        tokenServiceProto.NewTokenServiceClient(client),
		PermissionServiceClient:   permissionServiceProto.NewPermissionServiceClient(client),
		OrganizationServiceClient: organizationServiceProto.NewOrganizationServiceClient(client),
		InvitationServiceClient:   invitationServiceProto.NewInvitationServiceClient(client),
	}, teardown
}
