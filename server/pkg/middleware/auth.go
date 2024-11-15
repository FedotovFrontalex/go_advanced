package middleware

import (
	"context"
	"net/http"
	"server/configs"
	"server/pkg/jwt"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Autorization")

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

		ctx := context.WithValue(req.Context(), ContextEmailKey, data.Email)
		reqWithCtx := req.WithContext(ctx)

		next.ServeHTTP(w, reqWithCtx)
	})
}
