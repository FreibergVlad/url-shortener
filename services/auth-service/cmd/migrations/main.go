package main

import (
	"github.com/FreibergVlad/url-shortener/auth-service/internal/config"
	"github.com/FreibergVlad/url-shortener/auth-service/internal/db/migrations"
	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/rs/zerolog"
)

func main() {
	config := config.NewMigrationConfig()

	logLevel := must.Return(zerolog.ParseLevel(config.LogLevel))
	zerolog.SetGlobalLevel(logLevel)

	migrations.Migrate(config)
}
