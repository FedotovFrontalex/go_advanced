package configs

import (
	"errors"
	"orderApi/pkg/logger"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Db              DbConfig
	SessionIdLength int
	Auth            AuthConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		logger.Error(errors.New("failed to load env file. Using default config"))
	}

	sessionIdLength, err := strconv.ParseInt(os.Getenv("SESSION_ID_LENGTH"), 10, 32)

	logger.Message(sessionIdLength)

	if err != nil {
		logger.Error(errors.New("failed to parse session_id length from env. Use 20 by default"))
		sessionIdLength = 20
	}

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
		SessionIdLength: int(sessionIdLength),
	}
}
