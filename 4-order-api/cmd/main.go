package main

import (
	"net/http"
	"orderApi/configs"
	"orderApi/internal/auth"
	"orderApi/internal/order"
	"orderApi/internal/product"
	"orderApi/internal/user"
	"orderApi/pkg/db"
	"orderApi/pkg/logger"
	"orderApi/pkg/middleware"
)

func App() http.Handler {
	conf := configs.Load()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	logger.Message("initialize repositories")
	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)
	orderRepository := order.NewOrderRepository(database)

	logger.Message("initialize services")
	userService := user.NewUserService(&user.UserServiceDeps{
		Config:         conf,
		UserRepository: userRepository,
	})
	authService := auth.NewAuthService(&auth.AuthServiceDeps{
		Config:      conf,
		UserService: userService,
	})
	productService := product.NewProductService(productRepository)
	orderService := order.NewOrderService(&order.OrderServiceDeps{
		OrderRepository: orderRepository,
	})

	logger.Message("initialize routes")
	product.NewProductHandler(
		router,
		product.ProductHandlerDeps{
			Config:         conf,
			ProductService: productService,
		},
	)

	auth.NewAuthHandler(
		router,
		auth.AuthHandlerDeps{
			Config:      conf,
			AuthService: authService,
		},
	)

	order.NewOrderHandler(
		router,
		order.OrderHandlerDeps{
			Config:         conf,
			OrderService:   orderService,
			UserService:    userService,
			ProductService: productService,
		},
	)

	stack := middleware.Chain(
		middleware.InitCors(conf.Security),
		middleware.Log,
	)

	return stack(router)
}

func main() {
	logger.Message("initialize server start")
	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	logger.Message("Starting server on port 8081")
	err := server.ListenAndServe()

	if err != nil {
		logger.Error(err)
	}
}
