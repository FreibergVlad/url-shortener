package config

import "github.com/FreibergVlad/url-shortener/shared/go/pkg/config"

type PostgresConfig struct {
	DSN string `env:"DSN,notEmpty"`
}

type RedisConfig struct {
	DSN string `env:"DSN,notEmpty"`
}

type JWTConfig struct {
	Secret                 string `env:"SECRET,notEmpty"`
	RefreshSecret          string `env:"REFRESH_SECRET,notEmpty"`
	LifetimeSeconds        int    `env:"LIFETIME_SECONDS,notEmpty"`
	RefreshLifetimeSeconds int    `env:"REFRESH_LIFETIME_SECONDS,notEmpty"`
}

type AdminConfig struct {
	Email    string `env:"EMAIL,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
}

type IdentityServiceConfig struct {
	Port     int            `env:"PORT,notEmpty"`
	LogLevel string         `env:"LOG_LEVEL" envDefault:"error"`
	JWT      JWTConfig      `envPrefix:"JWT_"`
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	Admin    AdminConfig    `envPrefix:"ADMIN_"`
}

type MigrationConfig struct {
	LogLevel string         `env:"LOG_LEVEL" envDefault:"error"`
	Postgres PostgresConfig `envPrefix:"MIGRATION_POSTGRES_"`
}

func NewMigrationConfig() MigrationConfig {
	cfg := MigrationConfig{}
	config.ParseConfig(&cfg)
	return cfg
}

func NewIdentityServiceConfig() IdentityServiceConfig {
	cfg := IdentityServiceConfig{}
	config.ParseConfig(&cfg)
	return cfg
}
