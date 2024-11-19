package main

import (
	"orderApi/internal/order"
	"orderApi/internal/product"
	"orderApi/internal/user"
	"orderApi/pkg/logger"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger.Message("Migration start")
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&product.Product{},
		&user.User{},
		&order.Order{},
	)

	logger.Success("Migration end")
}
