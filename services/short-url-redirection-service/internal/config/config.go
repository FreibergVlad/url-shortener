package config

import "github.com/FreibergVlad/url-shortener/shared/go/pkg/config"

type Config struct {
	Port                         int    `env:"PORT,notEmpty"`
	LogLevel                     string `env:"LOG_LEVEL" envDefault:"error"`
	Domain                       string `env:"DOMAIN,notEmpty"`
	ShortURLManagementServiceDSN string `env:"SHORT_URL_MANAGEMENT_SERVICE_DSN,notEmpty"`
}

func New() Config {
	cfg := Config{}
	config.ParseConfig(&cfg)
	return cfg
}
