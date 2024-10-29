package main

import (
	"net/http"
	"server/configs"
	"server/internal/auth"
	"server/pkg/logger"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	logger.Message("Starting server on 8081 port")
	err := server.ListenAndServe()

	if err != nil {
		logger.Error(err)
	}
}
