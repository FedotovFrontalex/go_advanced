package main

import (
	"net/http"
	"server/configs"
	"server/internal/auth"
	"server/internal/link"
	"server/internal/user"
	"server/pkg/db"
	"server/pkg/logger"
	"server/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)

	authService := auth.NewAuthService(userRepository)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		Config:         conf,
		LinkRepository: linkRepository,
	})

	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	logger.Message("Starting server on 8081 port")
	err := server.ListenAndServe()

	if err != nil {
		logger.Error(err)
	}
}
