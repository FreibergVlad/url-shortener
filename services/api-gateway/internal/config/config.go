package config

import "github.com/FreibergVlad/url-shortener/shared/go/pkg/config"

type Config struct {
	Port                         int    `env:"PORT,notEmpty"`
	LogLevel                     string `env:"LOG_LEVEL" envDefault:"error"`
	JWTSecret                    string `env:"JWT_SECRET,notEmpty"`
	AuthServiceDSN               string `env:"AUTH_SERVICE_DSN,notEmpty"`
	DomainServiceDSN             string `env:"DOMAIN_SERVICE_DSN,notEmpty"`
	ShortUrlManagementServiceDSN string `env:"SHORT_URL_MANAGEMENT_SERVICE_DSN,notEmpty"`
	ShortUrlGeneratorServiceDSN  string `env:"SHORT_URL_GENERATOR_SERVICE_DSN,notEmpty"`
}

func New() Config {
	cfg := Config{}
	config.ParseConfig(&cfg)
	return cfg
}
