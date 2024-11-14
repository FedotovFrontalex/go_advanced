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
	Security        SecurityConfig
}

type DbConfig struct {
	Dsn string
}

type AuthConfig struct {
	Secret string
}

type SecurityConfig struct {
	Domains string
}

func Load() *Config {
	err := godotenv.Load()

	if err != nil {
		logger.Error(errors.New("failed to load env file. Using default config"))
	}

	sessionIdLength, err := strconv.ParseInt(os.Getenv("SESSION_ID_LENGTH"), 10, 32)

	logger.Message(sessionIdLength)

	if err != nil {
		logger.Error(errors.New(ErrNoSessionIdLength))
		sessionIdLength = 20
	}

	dsn := os.Getenv("DSN")

	if dsn == "" {
		panic(ErrNoDSN)
	}

	secret := os.Getenv("SECRET")

	if secret == "" {
		logger.Error(errors.New(ErrNoSecret))
	}

	domains := os.Getenv("DOMAINS")

	if domains == "" {
		logger.Error(errors.New(ErrNoDomains))
	}

	return &Config{
		Db: DbConfig{
			Dsn: dsn,
		},
		Auth: AuthConfig{
			Secret: secret,
		},
		Security: SecurityConfig{
			Domains: domains,
		},
		SessionIdLength: int(sessionIdLength),
	}
}
