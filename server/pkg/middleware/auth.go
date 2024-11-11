package middleware

import (
	"net/http"
	"server/pkg/logger"
	"strings"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Autorization")
		authToken := strings.TrimPrefix(authHeader, "Bearer ")
		logger.Log(authToken)
		next.ServeHTTP(w, req)
	})
}
