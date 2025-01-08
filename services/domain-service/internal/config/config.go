package config

import "github.com/FreibergVlad/url-shortener/shared/go/pkg/config"

type Config struct {
	Port           int      `env:"PORT,notEmpty"`
	LogLevel       string   `env:"LOG_LEVEL" envDefault:"error"`
	AuthServiceDSN string   `env:"AUTH_SERVICE_DSN,notEmpty"`
	PublicDomains  []string `env:"PUBLIC_DOMAINS,notEmpty"`
}

func New() Config {
	cfg := Config{}
	config.ParseConfig(&cfg)
	return cfg
}
