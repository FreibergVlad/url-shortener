package invitations

import (
	"context"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories/invitations/gen"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/schema"
)

type Repository interface {
	Create(ctx context.Context, invitation *schema.Invitation) error
	GetByID(ctx context.Context, id string) (*schema.Invitation, error)
	UpdateStatusByID(ctx context.Context, id, status string) error
}

type PostgresRepository struct {
	db gen.DBTX
}

func NewInvitationRepository(db gen.DBTX) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, invitation *schema.Invitation) error {
	querier := gen.New(r.db)
	params := gen.CreateInvitationParams{
		ID:             repositories.MustPgUUIDFromString(invitation.ID),
		Status:         invitation.Status,
		OrganizationID: repositories.MustPgUUIDFromString(invitation.OrganizationID),
		Email:          invitation.Email,
		RoleID:         invitation.RoleID,
		CreatedAt:      repositories.MustPgTimestamptzFromTime(invitation.CreatedAt),
		ExpiresAt:      repositories.MustPgTimestamptzFromTime(invitation.ExpiresAt),
		CreatedBy:      repositories.MustPgUUIDFromString(invitation.CreatedBy),
	}

	return repositories.TranslatePgError(querier.CreateInvitation(ctx, params))
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*schema.Invitation, error) {
	querier := gen.New(r.db)

	row, err := querier.GetInvitationById(ctx, repositories.MustPgUUIDFromString(id))
	if err != nil {
		return nil, repositories.TranslatePgError(err)
	}

	return &schema.Invitation{
		ID:             repositories.MustPgUUIDToString(row.ID),
		Status:         row.Status,
		OrganizationID: repositories.MustPgUUIDToString(row.OrganizationID),
		Email:          row.Email,
		RoleID:         row.RoleID,
		CreatedAt:      row.CreatedAt.Time,
		ExpiresAt:      row.ExpiresAt.Time,
		CreatedBy:      repositories.MustPgUUIDToString(row.CreatedBy),
	}, nil
}

func (r *PostgresRepository) UpdateStatusByID(ctx context.Context, id, status string) error {
	querier := gen.New(r.db)
	params := gen.UpdateInvitationStatusByIdParams{
		ID:     repositories.MustPgUUIDFromString(id),
		Status: status,
	}

	return repositories.TranslatePgError(querier.UpdateInvitationStatusById(ctx, params))
}
