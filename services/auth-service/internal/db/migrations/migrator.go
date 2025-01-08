package migrations

import (
	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

const migrationsPath = "file:///db/migrations/versions/"

func Migrate(config config.MigrationConfig) {
	migrator := must.Return(migrate.New(migrationsPath, config.Postgres.DSN))

	err := migrator.Up()
	if err == nil {
		log.Info().Msg("migrations applied successfully")
		return
	}
	if err == migrate.ErrNoChange {
		log.Info().Msg("no new migrations to apply")
		return
	}
	log.Fatal().Err(err).Msg("failed to apply migrations")
}
