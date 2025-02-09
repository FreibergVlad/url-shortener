package testing

import (
	"context"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/stretchr/testify/mock"
)

type MockedUserRepository struct {
	mock.Mock
}

func (r *MockedUserRepository) Create(ctx context.Context, user *schema.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *MockedUserRepository) GetByID(ctx context.Context, id string) (*schema.User, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(*schema.User), args.Error(1)
}

func (r *MockedUserRepository) GetByEmail(ctx context.Context, email string) (*schema.User, error) {
	args := r.Called(ctx, email)
	return args.Get(0).(*schema.User), args.Error(1)
}

type MockedOrganizationRepository struct {
	mock.Mock
}

func (r *MockedOrganizationRepository) Create(ctx context.Context, organization *schema.Organization) error {
	args := r.Called(ctx, organization)
	return args.Error(0)
}

func (r *MockedOrganizationRepository) CreateOrganizationMembership(
	ctx context.Context, params *organizations.CreateOrganizationMembershipParams,
) error {
	args := r.Called(ctx, params)
	return args.Error(0)
}

func (r *MockedOrganizationRepository) ListOrganizationMembershipsByUserID(
	ctx context.Context, userID string,
) (schema.OrganizationMemberships, error) {
	args := r.Called(ctx, userID)
	return args.Get(0).(schema.OrganizationMemberships), args.Error(1)
}
