package repositories

import (
	"fmt"
	"time"

	serviceErrors "github.com/FreibergVlad/url-shortener/shared/go/pkg/errors"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func MustPgUUIDFromString(uuid string) pgtype.UUID {
	var pgUUID pgtype.UUID
	must.Do(pgUUID.Scan(uuid))
	return pgUUID
}

func MustPgUUIDToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}

func MustPgTimestamptzFromTime(time time.Time) pgtype.Timestamptz {
	var pgTmz pgtype.Timestamptz
	must.Do(pgTmz.Scan(time))
	return pgTmz
}

func TranslatePgError(err error) error {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		if err == pgx.ErrNoRows {
			return serviceErrors.ErrResourceNotFound
		}
		return err
	}
	if pgErr.Code == pgerrcode.UniqueViolation {
		return serviceErrors.ErrDuplicateResource
	}
	return err
}
