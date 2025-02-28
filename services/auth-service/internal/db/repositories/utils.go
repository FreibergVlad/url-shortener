package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
)

func MustPgUUIDFromString(uuid string) pgtype.UUID {
	var pgUUID pgtype.UUID
	must.Do(pgUUID.Scan(uuid))
	return pgUUID
}

func MustPgUUIDToString(uuid pgtype.UUID) string {
	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		uuid.Bytes[0:4],
		uuid.Bytes[4:6],
		uuid.Bytes[6:8],
		uuid.Bytes[8:10],
		uuid.Bytes[10:16],
	)
}

func MustPgTimestamptzFromTime(time time.Time) pgtype.Timestamptz {
	var pgTmz pgtype.Timestamptz
	must.Do(pgTmz.Scan(time))
	return pgTmz
}

func TranslatePgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if pgErr.Code == pgerrcode.UniqueViolation {
		return ErrAlreadyExists
	}
	return pgErr
}
