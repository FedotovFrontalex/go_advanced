package auth

import (
	"net/http"
	"orderApi/configs"
	"orderApi/internal/user"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/jwt"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
)

type AuthServiceInterface interface {
	CreateSession(string) (string, error)
	VerifySession(string, int) (*user.User, *apierrors.Error)
}

type AuthHandlerDeps struct {
	*configs.Config
	AuthService AuthServiceInterface
}

type AuthHandler struct {
	*configs.Config
	AuthService AuthServiceInterface
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	router.HandleFunc("POST /auth", handler.Auth())
	router.HandleFunc("POST /auth/verify", handler.VerifyAuth())
}

func (handler *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[AuthRequest](&w, req)

		if err != nil {
			return
		}

		sessionId, err := handler.AuthService.CreateSession(body.Phone)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData := AuthResponse{
			SessionId: sessionId,
		}

		response.Json(w, responseData, 200)
	}
}

func (handler *AuthHandler) VerifyAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[AuthVerifyRequest](&w, req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, apierr := handler.AuthService.VerifySession(body.SessionId, body.Code)

		if apierr != nil {
			http.Error(w, apierr.Error(), apierr.GetStatus())
			return
		}

		jwtService := jwt.NewJWT(handler.Config.Auth.Secret)
		token, err := jwtService.Create(jwt.JWTData{
			Phone:     user.Phone,
			SessionId: user.SessionId,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		responseData := &AuthVerificationResponse{
			Token: token,
		}

		response.Json(w, responseData, http.StatusOK)
	}
}
