package config

import "github.com/FreibergVlad/url-shortener/shared/go/pkg/config"

type MongoDBConfig struct {
	DSN    string `env:"DSN,notEmpty"`
	DBName string `env:"DB_NAME,notEmpty"`
}

type RedisConfig struct {
	DSN string `env:"DSN,notEmpty"`
}

type ShortURLGeneratorConfig struct {
	Port                  int           `env:"PORT,notEmpty"`
	LogLevel              string        `env:"LOG_LEVEL" envDefault:"error"`
	ShortURLAliasLength   int           `env:"SHORT_URL_ALIAS_LENGTH" envDefault:"4"`
	ShortURLScheme        string        `env:"SHORT_URL_SCHEME" envDefault:"https"`
	MaxRetriesOnCollision int           `env:"MAX_RETRIES_ON_COLLISION" envDefault:"5"`
	MongoDB               MongoDBConfig `envPrefix:"MONGO_"`
	Redis                 RedisConfig   `envPrefix:"REDIS_"`
	AuthServiceDSN        string        `env:"AUTH_SERVICE_DSN,notEmpty"`
	DomainServiceDSN      string        `env:"DOMAIN_SERVICE_DSN,notEmpty"`
}

type ShortURLManagerConfig struct {
	Port           int           `env:"PORT,notEmpty"`
	LogLevel       string        `env:"LOG_LEVEL" envDefault:"error"`
	MongoDB        MongoDBConfig `envPrefix:"MONGO_"`
	Redis          RedisConfig   `envPrefix:"REDIS_"`
	AuthServiceDSN string        `env:"AUTH_SERVICE_DSN,notEmpty"`
}

type MigrationConfig struct {
	LogLevel string        `env:"LOG_LEVEL" envDefault:"error"`
	MongoDB  MongoDBConfig `envPrefix:"MONGO_"`
}

func NewMigrationConfig() MigrationConfig {
	cfg := MigrationConfig{}
	config.ParseConfig(&cfg)
	return cfg
}

func NewShortURLManagerConfig() ShortURLManagerConfig {
	cfg := ShortURLManagerConfig{}
	config.ParseConfig(&cfg)
	return cfg
}

func NewShortURLGeneratorConfig() ShortURLGeneratorConfig {
	cfg := ShortURLGeneratorConfig{}
	config.ParseConfig(&cfg)
	return cfg
}
