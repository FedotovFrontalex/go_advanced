package auth

import (
	"net/http"
	"os"
	"server/configs"
	"server/pkg/logger"
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
		result := LoginResponse{
			Token: os.Getenv("TOKEN"),
		}
		
		response.Json(w, result, 200)
	}
}

func (handler *AuthHandler) register() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("REGISTER"))
		logger.Message("register")
	}
}
