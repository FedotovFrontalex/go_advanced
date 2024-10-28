package main

import (
	"net/http"
	"simpleServer/logger"
)

func main() {
	router := http.NewServeMux()
	NewHelloHandler(router)	

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	logger.Message("Server is listening on port 8081")
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err)
	}
}
