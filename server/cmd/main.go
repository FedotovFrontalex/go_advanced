package main

import (
	"net/http"
	"server/configs"
	"server/internal/auth"
	"server/internal/link"
	"server/pkg/db"
	"server/pkg/logger"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	linkRepository := link.NewLinkRepository(database)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		Config:         conf,
		LinkRepository: linkRepository,
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
