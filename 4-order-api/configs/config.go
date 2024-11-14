package configs

import (
	"errors"
	"orderApi/pkg/logger"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db DbConfig
}

type DbConfig struct {
	Dsn string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		logger.Error(errors.New("Failed to load env file. Using default config"))
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
	}
}
