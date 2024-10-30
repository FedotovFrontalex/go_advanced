package main

import (
	"net/http"
	"validationApi/configs"
	"validationApi/internal/verify"
	"validationApi/pkg/logger"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config: conf,
	})
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	logger.Message("Starting on port 8081")
	err := server.ListenAndServe()
	logger.Error(err)
}
