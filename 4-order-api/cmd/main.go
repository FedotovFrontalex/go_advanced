package main

import (
	"net/http"
	"orderApi/configs"
	"orderApi/internal/product"
	"orderApi/pkg/db"
	"orderApi/pkg/logger"
	"orderApi/pkg/middleware"
)

func main() {
	logger.Message("initialize server start")
	conf := configs.Load()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	logger.Message("initialize repositories")
	productRepository := product.NewProductRepository(database)

	logger.Message("initialize routes")
	product.NewProductHandler(
		router,
		&product.ProductHandlerDeps{
			Config:            conf,
			ProductRepository: productRepository,
		},
	)

	middlewareChain := middleware.Chain(
		middleware.CORS,
		middleware.Log,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: middlewareChain(router),
	}

	logger.Message("Starting server on port 8081")
	err := server.ListenAndServe()

	if err != nil {
		logger.Error(err)
	}
}
