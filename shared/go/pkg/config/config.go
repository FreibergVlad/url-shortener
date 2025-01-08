package config

import (
	"os"

	"github.com/FreibergVlad/url-shortener/shared/go/pkg/must"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// ParseConfig loads environment variables into the provided config structure.
// It attempts to load environment variables from a .env file if specified by the
// ENV_FILE environment variable, falling back to system environment variables
// if the file is missing or fails to load. If parsing the environment variables
// into the config structure fails, the function will panic.
func ParseConfig(config interface{}) {
	envFile := os.Getenv("ENV_FILE")
	// ignore error loading .env file, fallback to env vars instead
	godotenv.Load(envFile) //nolint:errcheck
	must.Do(env.Parse(config))
}
