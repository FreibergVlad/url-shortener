package migrations

import (
	"errors"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb" // load MongoDB driver.
	_ "github.com/golang-migrate/migrate/v4/source/file"      // load file source driver.
	"github.com/rs/zerolog/log"
)

const migrationsPath = "file:///db/migrations/versions/"

func Migrate(config config.MigrationConfig) {
	migrator := must.Return(migrate.New(migrationsPath, config.MongoDB.DSN))

	err := migrator.Up()
	if err == nil {
		log.Info().Msg("Migrations applied successfully")
		return
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Info().Msg("No new migrations to apply")
		return
	}

	log.Fatal().Err(err).Msg("Failed to apply migrations")
}
