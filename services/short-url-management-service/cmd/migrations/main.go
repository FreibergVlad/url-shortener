package main

import (
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/config"
	"github.com/FreibergVlad/url-shortener/short-url-management-service/internal/db/migrations"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
)

func main() {
	config := config.NewMigrationConfig()

	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	migrations.Migrate(config)
}
