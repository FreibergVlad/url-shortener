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
	"github.com/rs/zerolog/log"
)

type Repository interface {
	Create(ctx context.Context, organization *schema.Organization) error
	CreateOrganizationMembership(ctx context.Context, params *CreateOrganizationMembershipParams) error
	ListOrganizationMembershipsByUserID(ctx context.Context, userID string) (schema.OrganizationMemberships, error)
}

type CreateOrganizationMembershipParams struct {
	OrganizationID string
	UserID         string
	RoleID         string
	CreatedAt      time.Time
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, organization *schema.Organization) error {
	var err error

	transactionCtx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return repositories.TranslatePgError(err)
	}

	defer func() {
		if transactionCtx == nil {
			return
		}

		if err != nil {
			if rollbackErr := transactionCtx.Rollback(ctx); rollbackErr != nil {
				log.Err(rollbackErr).Msg("failed to rollback transaction")
			}
		} else {
			if commitErr := transactionCtx.Commit(ctx); commitErr != nil {
				log.Err(commitErr).Msg("failed to rollback transaction")
				err = commitErr
			}
		}
	}()

	querier := gen.New(transactionCtx)

	organizationID := repositories.MustPgUUIDFromString(organization.ID)
	createdBy := repositories.MustPgUUIDFromString(organization.CreatedBy)
	createdAt := repositories.MustPgTimestamptzFromTime(organization.CreatedAt)

	createOrgParams := gen.CreateOrganizationParams{
		ID:        organizationID,
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
		OrganizationID: organizationID,
		UserID:         createdBy,
		RoleID:         roles.RoleIDOwner,
		CreatedAt:      createdAt,
	}
	err = querier.CreateOrganizationMembership(ctx, createOrgMembershipParams)
	if err != nil {
		return repositories.TranslatePgError(err)
	}

	return nil
}

func (r *PostgresRepository) CreateOrganizationMembership(
	ctx context.Context, params *CreateOrganizationMembershipParams,
) error {
	querier := gen.New(r.db)
	createOrgMembershipParams := gen.CreateOrganizationMembershipParams{
		OrganizationID: repositories.MustPgUUIDFromString(params.OrganizationID),
		UserID:         repositories.MustPgUUIDFromString(params.UserID),
		RoleID:         params.RoleID,
		CreatedAt:      repositories.MustPgTimestamptzFromTime(params.CreatedAt),
	}

	return repositories.TranslatePgError(querier.CreateOrganizationMembership(ctx, createOrgMembershipParams))
}

func (r *PostgresRepository) ListOrganizationMembershipsByUserID(
	ctx context.Context, userID string,
) (schema.OrganizationMemberships, error) {
	querier := gen.New(r.db)

	rows, err := querier.ListOrganizationMembershipsByUserId(ctx, repositories.MustPgUUIDFromString(userID))
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
