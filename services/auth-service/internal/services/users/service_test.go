package users_test

import (
	"context"
	"testing"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/users"
	testUtils "github.com/FreibergVlad/url-shortener/auth-service/internal/testing"
	userServiceMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()
	userRepo := testUtils.MockedUserRepository{}
	userService := users.New(&userRepo, clock.NewFixedClock(fixedTime))
	ctx := context.Background()
	request := testUtils.CreateTestUserRequest()

	userRepo.On("Create", ctx, mock.Anything).Return(nil)

	response, err := userService.CreateUser(ctx, request)

	require.NoError(t, err)

	assert.Equal(t, request.Email, response.User.Email)
	assert.Equal(t, request.FullName, response.User.FullName)
	assert.Equal(t, timestamppb.New(fixedTime), response.User.CreatedAt)
	assert.Equal(t, roles.RoleProvisional.ID, response.User.Role.Id)
	assert.Equal(t, roles.RoleProvisional.Description, response.User.Role.Description)

	userRepo.AssertExpectations(t)
}

func TestCreateUserWhenPasswordHashingError(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	userService := users.New(&userRepo, clock.NewFixedClock(time.Now()))
	ctx := context.Background()
	request := userServiceMessages.CreateUserRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 80),
	}

	response, err := userService.CreateUser(ctx, &request)

	require.ErrorIs(t, err, bcrypt.ErrPasswordTooLong)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
}

func TestCreateUserWhenDatabaseError(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	userService := users.New(&userRepo, clock.NewFixedClock(time.Now()))
	ctx := context.Background()
	wantErr := gofakeit.ErrorDatabase()

	userRepo.On("Create", ctx, mock.Anything).Return(wantErr)

	response, err := userService.CreateUser(ctx, testUtils.CreateTestUserRequest())

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
}

func TestGetMe(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	userService := users.New(&userRepo, clock.NewFixedClock(time.Now()))
	user := schema.User{ID: gofakeit.UUID(), Email: gofakeit.Email(), FullName: gofakeit.Name()}
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), user.ID)

	userRepo.On("GetByID", ctx, user.ID).Return(&user, nil)

	response, err := userService.GetMe(ctx, &userServiceMessages.GetMeRequest{})

	require.NoError(t, err)

	assert.Equal(t, user.ID, response.User.Id)
	assert.Equal(t, user.Email, response.User.Email)
	assert.Equal(t, user.FullName, response.User.FullName)

	userRepo.AssertExpectations(t)
}

func TestGetMeWhenDatabaseError(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	userService := users.New(&userRepo, clock.NewFixedClock(time.Now()))
	userID := gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)
	wantErr := gofakeit.ErrorDatabase()

	userRepo.On("GetByID", ctx, userID).Return(&schema.User{}, wantErr)

	response, err := userService.GetMe(ctx, &userServiceMessages.GetMeRequest{})

	require.ErrorIs(t, err, wantErr)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
}

func TestGetMeWhenUserNotFound(t *testing.T) {
	t.Parallel()

	userRepo := testUtils.MockedUserRepository{}
	userService := users.New(&userRepo, clock.NewFixedClock(time.Now()))
	userID := gofakeit.UUID()
	ctx := grpcUtils.IncomingContextWithUserID(context.Background(), userID)

	userRepo.On("GetByID", ctx, userID).Return(&schema.User{}, repositories.ErrNotFound)

	response, err := userService.GetMe(ctx, &userServiceMessages.GetMeRequest{})

	require.ErrorIs(t, err, errors.ErrUserNotFound)

	assert.Nil(t, response)

	userRepo.AssertExpectations(t)
}
