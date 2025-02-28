package repositories_test

import (
	"testing"

	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/repositories"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestTranslatePgError(t *testing.T) {
	t.Parallel()

	unknownErr := gofakeit.Error()
	unknownPgconnErr := &pgconn.PgError{Code: pgerrcode.InternalError}

	tests := []struct {
		name    string
		gotErr  error
		wantErr error
	}{
		{name: "not a pgconn.PgError or pgx instance", gotErr: unknownErr, wantErr: unknownErr},
		{name: "pgx.ErrNoRows", gotErr: pgx.ErrNoRows, wantErr: repositories.ErrNotFound},
		{
			name:    "pgconn.PgError (code = UniqueViolation)",
			gotErr:  &pgconn.PgError{Code: pgerrcode.UniqueViolation},
			wantErr: repositories.ErrAlreadyExists,
		},
		{name: "pgconn.PgError (code = InternalError)", gotErr: unknownPgconnErr, wantErr: unknownPgconnErr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.ErrorIs(t, repositories.TranslatePgError(test.gotErr), test.wantErr)
		})
	}
}
