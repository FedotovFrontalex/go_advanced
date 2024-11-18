package main

import (
	"net/http"
	"server/configs"
	"server/internal/auth"
	"server/internal/link"
	"server/internal/stat"
	"server/internal/user"
	"server/pkg/db"
	"server/pkg/event"
	"server/pkg/logger"
	"server/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		Config:         conf,
		LinkRepository: linkRepository,
		EventBus:       eventBus,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{
		Config:         conf,
		StatRepository: statRepository,
	})

	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	go statService.AddClick()

	logger.Message("Starting server on 8081 port")
	err := server.ListenAndServe()

	if err != nil {
		logger.Error(err)
	}
}
