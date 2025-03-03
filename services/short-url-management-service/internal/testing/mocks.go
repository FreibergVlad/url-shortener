package testing

import (
	"context"

	domainServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/domains/messages/v1"
	userServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockedDomainServiceClient struct {
	mock.Mock
}

func NewMockedDomainServiceClient(fqdns []string) *MockedDomainServiceClient {
	domains := make([]*domainServiceMessages.Domain, len(fqdns))
	for i, fqdn := range fqdns {
		domains[i] = &domainServiceMessages.Domain{Fqdn: fqdn}
	}

	mockedResponse := &domainServiceMessages.ListOrganizationDomainResponse{Data: domains, Total: int64(len(domains))}

	c := &MockedDomainServiceClient{}
	c.On("ListOrganizationDomain", mock.Anything, mock.Anything, mock.Anything).Return(mockedResponse, nil)

	return c
}

func (c *MockedDomainServiceClient) ListOrganizationDomain(
	ctx context.Context, req *domainServiceMessages.ListOrganizationDomainRequest, opts ...grpc.CallOption,
) (*domainServiceMessages.ListOrganizationDomainResponse, error) {
	args := c.Called(ctx, req, opts)
	return args.Get(0).(*domainServiceMessages.ListOrganizationDomainResponse), args.Error(1)
}

type MockedUserServiceClient struct {
	mock.Mock
}

func NewMockedUserServiceClient(user schema.User) *MockedUserServiceClient {
	response := &userServiceMessages.GetMeResponse{User: &userServiceMessages.User{Id: user.ID, Email: user.Email}}

	c := &MockedUserServiceClient{}
	c.On("GetMe", mock.Anything, mock.Anything, mock.Anything).Return(response, nil)

	return c
}

func (c *MockedUserServiceClient) GetMe(
	ctx context.Context, req *userServiceMessages.GetMeRequest, opts ...grpc.CallOption,
) (*userServiceMessages.GetMeResponse, error) {
	args := c.Called(ctx, req, opts)
	return args.Get(0).(*userServiceMessages.GetMeResponse), args.Error(1)
}
