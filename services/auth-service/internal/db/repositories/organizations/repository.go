package organizations

import (
	"context"
	"time"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/organizations/gen"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/services/roles"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, organization *schema.Organization) error
	CreateOrganizationMembership(ctx context.Context, params *CreateOrganizationMembershipParams) error
	ListOrganizationMembershipsByUserId(ctx context.Context, userId string) (schema.OrganizationMemberships, error)
}

type CreateOrganizationMembershipParams struct {
	OrganizationID string
	UserID         string
	RoleID         string
	CreatedAt      time.Time
}

type organizationRepository struct {
	db *pgxpool.Pool
}

func NewOrganizationRepository(db *pgxpool.Pool) *organizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, organization *schema.Organization) (err error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return repositories.TranslatePgError(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	querier := gen.New(tx)

	id := repositories.MustPgUUIDFromString(organization.ID)
	createdBy := repositories.MustPgUUIDFromString(organization.CreatedBy)
	createdAt := repositories.MustPgTimestamptzFromTime(organization.CreatedAt)

	createOrgParams := gen.CreateOrganizationParams{
		ID:        id,
		Name:      organization.Name,
		Slug:      organization.Slug,
		CreatedAt: createdAt,
		CreatedBy: createdBy,
	}
	err = querier.CreateOrganization(ctx, createOrgParams)
	if err != nil {
		return repositories.TranslatePgError(err)
	}

	createOrgMembershipParams := gen.CreateOrganizationMembershipParams{
		OrganizationID: id,
		UserID:         createdBy,
		RoleID:         roles.RoleIdOwner,
		CreatedAt:      createdAt,
	}
	err = querier.CreateOrganizationMembership(ctx, createOrgMembershipParams)
	if err != nil {
		return repositories.TranslatePgError(err)
	}

	return nil
}

func (r *organizationRepository) CreateOrganizationMembership(ctx context.Context, params *CreateOrganizationMembershipParams) error {
	querier := gen.New(r.db)
	createOrgMembershipParams := gen.CreateOrganizationMembershipParams{
		OrganizationID: repositories.MustPgUUIDFromString(params.OrganizationID),
		UserID:         repositories.MustPgUUIDFromString(params.UserID),
		RoleID:         params.RoleID,
		CreatedAt:      repositories.MustPgTimestamptzFromTime(params.CreatedAt),
	}

	return repositories.TranslatePgError(querier.CreateOrganizationMembership(ctx, createOrgMembershipParams))
}

func (r *organizationRepository) ListOrganizationMembershipsByUserId(ctx context.Context, userId string) (schema.OrganizationMemberships, error) {
	querier := gen.New(r.db)

	rows, err := querier.ListOrganizationMembershipsByUserId(ctx, repositories.MustPgUUIDFromString(userId))
	if err != nil {
		return nil, repositories.TranslatePgError(err)
	}

	memberships := schema.OrganizationMemberships{}
	for _, row := range rows {
		membership := schema.OrganizationMembership{
			Organization: schema.ShortOrganization{
				ID:   repositories.MustPgUUIDToString(row.OrganizationID),
				Slug: row.OrganizationSlug,
			},
			RoleID:    row.RoleID,
			CreatedAt: row.CreatedAt.Time,
		}
		memberships = append(memberships, &membership)
	}

	return memberships, nil
}
