package users

import (
	"context"

	userServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	"google.golang.org/grpc"
)

type UserInfoGetter interface {
	GetMe(
		ctx context.Context,
		req *userServiceMessages.GetMeRequest,
		opts ...grpc.CallOption,
	) (*userServiceMessages.GetMeResponse, error)
}

type GRPCServiceClient struct {
	grpcClient UserInfoGetter
}

func NewServiceClient(grpcClient UserInfoGetter) *GRPCServiceClient {
	return &GRPCServiceClient{
		grpcClient: grpcClient,
	}
}

func (c *GRPCServiceClient) GetUserInfo(ctx context.Context, userID string) (schema.User, error) {
	response, err := c.grpcClient.GetMe(
		grpcUtils.OutgoingContextWithUserID(ctx, userID),
		&userServiceMessages.GetMeRequest{},
	)
	if err != nil {
		return schema.User{}, err
	}

	return schema.User{
		ID:    response.User.Id,
		Email: response.User.Email,
	}, nil
}
