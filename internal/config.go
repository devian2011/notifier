package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"

	"notifications/internal/io/web"
)

type Config struct {
	Web                   web.Config
	StorageDsn            string `env:"APP_STORAGE_DSN" envDefault:"file://PWD/templates"`
	TransportsCfgFilePath string `env:"APP_TRANSPORTS_CFG_PATH" envDefault:"./config/transports.local.yaml"`
}

func loadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	parseErr := env.Parse(cfg)

	return cfg, parseErr
}
