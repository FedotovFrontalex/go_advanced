package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"orderApi/internal/product"
	"orderApi/pkg/logger"
	"os"
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

	db.AutoMigrate(&product.Product{})

	logger.Success("Migration end")
}
