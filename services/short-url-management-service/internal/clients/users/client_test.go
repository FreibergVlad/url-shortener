package users_test

import (
	"context"
	"testing"

	userServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/clients/users"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/schema"
	testUtils "github.com/FreibergVlad/url-shortener/short-url-management-service/internal/testing"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetUserInfo(t *testing.T) {
	t.Parallel()

	wantUser := schema.User{ID: gofakeit.UUID(), Email: gofakeit.Email()}
	grpcClient := testUtils.NewMockedUserServiceClient(wantUser)
	client := users.NewServiceClient(grpcClient)

	gotUser, err := client.GetUserInfo(context.Background(), wantUser.ID)

	require.NoError(t, err)
	require.Equal(t, wantUser, gotUser)
}

func TestGetUserInfoWhenGrpcError(t *testing.T) {
	t.Parallel()

	wantErr := gofakeit.ErrorGRPC()
	grpcClient := &testUtils.MockedUserServiceClient{}
	grpcClient.
		On("GetMe", mock.Anything, mock.Anything, mock.Anything).
		Return(&userServiceMessages.GetMeResponse{}, wantErr)
	client := users.NewServiceClient(grpcClient)

	_, err := client.GetUserInfo(context.Background(), gofakeit.UUID())

	require.ErrorIs(t, err, wantErr)
}
