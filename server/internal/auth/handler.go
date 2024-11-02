package auth

import (
	"net/http"
	"os"
	"server/configs"
	"server/pkg/logger"
	"server/pkg/request"
	"server/pkg/response"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.login())
	router.HandleFunc("POST /auth/register", handler.register())
}

func (handler *AuthHandler) login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("login")

		_, err := request.HandleBody[LoginRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		result := LoginResponse{
			Token: os.Getenv("TOKEN"),
		}

		response.Json(w, result, 200)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("register")

		_, err := request.HandleBody[RegisterRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		w.Write([]byte("REGISTER SUCCESS"))
	}
}
