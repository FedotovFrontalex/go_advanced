package main

import (
	"net/http"
	"randomApi/logger"
)

func main() {
	router := http.NewServeMux()
	NewApiHandler(router)
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	logger.Message("Server starting on 8081 port...")
	err := server.ListenAndServe()

	if err != nil {
		logger.Error(err)
	}
}
