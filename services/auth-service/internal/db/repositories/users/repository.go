package users

import (
	"context"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/users/gen"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
)

type Repository interface {
	Create(ctx context.Context, user *schema.User) error
	GetById(ctx context.Context, id string) (*schema.User, error)
	GetByEmail(ctx context.Context, email string) (*schema.User, error)
}

type userRepository struct {
	db gen.DBTX
}

func NewUserRepository(db gen.DBTX) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *schema.User) error {
	querier := gen.New(r.db)

	params := gen.CreateUserParams{
		ID:           repositories.MustPgUUIDFromString(user.ID),
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		RoleID:       user.RoleID,
		CreatedAt:    repositories.MustPgTimestamptzFromTime(user.CreatedAt),
	}

	return repositories.TranslatePgError(querier.CreateUser(ctx, params))
}

func (r *userRepository) GetById(ctx context.Context, id string) (*schema.User, error) {
	querier := gen.New(r.db)

	user, err := querier.GetUserById(ctx, repositories.MustPgUUIDFromString(id))
	if err != nil {
		return nil, repositories.TranslatePgError(err)
	}

	return r.userFromRow(user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*schema.User, error) {
	querier := gen.New(r.db)

	user, err := querier.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, repositories.TranslatePgError(err)
	}

	return r.userFromRow(user), nil
}

func (r *userRepository) userFromRow(user gen.User) *schema.User {
	return &schema.User{
		ID:           repositories.MustPgUUIDToString(user.ID),
		PasswordHash: user.PasswordHash,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		RoleID:       user.RoleID,
		CreatedAt:    user.CreatedAt.Time,
	}
}
