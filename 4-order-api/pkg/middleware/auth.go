package middleware

import (
	"context"
	"net/http"
	"orderApi/configs"
	"orderApi/pkg/jwt"
	"strings"
)

type key string

const (
	ContextSessionIdKey key = "ContextSessionId"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthed(w)
			return
		}

		authToken := strings.TrimPrefix(authHeader, "Bearer ")

		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(authToken)

		if !isValid {
			writeUnauthed(w)
			return
		}

		ctx := context.WithValue(req.Context(), ContextSessionIdKey, data.SessionId)
		reqWithCtx := req.WithContext(ctx)

		next.ServeHTTP(w, reqWithCtx)
	})
}
