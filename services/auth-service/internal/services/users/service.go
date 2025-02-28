package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	protoMessages "github.com/FreibergVlad/url-shortener/proto/pkg/users/messages/v1"
	protoService "github.com/FreibergVlad/url-shortener/proto/pkg/users/service/v1"
	grpcUtils "github.com/FreibergVlad/url-shortener/shared/go/pkg/api/grpc/utils"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/clock"
	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const bcryptCost = 10

type UserService struct {
	protoService.UnimplementedUserServiceServer
	userRepository users.Repository
	clock          clock.Clock
}

func New(userRepository users.Repository, clock clock.Clock) *UserService {
	return &UserService{userRepository: userRepository, clock: clock}
}

func (s *UserService) CreateUser(
	ctx context.Context, req *protoMessages.CreateUserRequest,
) (*protoMessages.CreateUserResponse, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := schema.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(passHash),
		RoleID:       roles.RoleIDProvisional,
		CreatedAt:    s.clock.Now(),
	}

	err = s.userRepository.Create(ctx, &user)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return nil, serviceErrors.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return &protoMessages.CreateUserResponse{User: userToProto(&user)}, nil
}

func (s *UserService) GetMe(ctx context.Context, _ *protoMessages.GetMeRequest) (*protoMessages.GetMeResponse, error) {
	userID := grpcUtils.UserIDFromIncomingContext(ctx)

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, serviceErrors.ErrUserNotFound
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return &protoMessages.GetMeResponse{User: userToProto(user)}, nil
}

func userToProto(user *schema.User) *protoMessages.User {
	return &protoMessages.User{
		Id:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      roles.GetRoleProto(user.RoleID),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
