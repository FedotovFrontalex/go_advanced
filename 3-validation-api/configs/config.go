package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
		Email string
		Password string
		Address string
}

func LoadConfig() *Config {
		godotenv.Load()

		return &Config{
				Email: os.Getenv("EMAIL"),
				Password: os.Getenv("PASSWORD"),
				Address: os.Getenv("ADDRESS"),
		}
}
