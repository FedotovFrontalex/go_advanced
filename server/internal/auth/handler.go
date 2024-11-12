package auth

import (
	"net/http"
	"server/configs"
	"server/pkg/jwt"
	"server/pkg/logger"
	"server/pkg/request"
	"server/pkg/response"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("login")

		body, err := request.HandleBody[LoginRequest](&w, req)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email, err := handler.AuthService.Login(body.Email, body.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jwtService := jwt.NewJWT(handler.Auth.Secret)
		token, err := jwtService.Create(email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData := LoginResponse{
			Token: token,
		}

		response.Json(w, responseData, http.StatusOK)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("register")

		body, err := request.HandleBody[RegisterRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jwtService := jwt.NewJWT(handler.Config.Auth.Secret)
		token, err := jwtService.Create(email)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData := &RegisterResponse{
			Token: token,
		}

		response.Json(w, responseData, http.StatusOK)
	}
}
