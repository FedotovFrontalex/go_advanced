package middleware

import (
	"net/http"
	"server/pkg/logger"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, req)
		logger.Log(wrapper.StatusCode, req.Method, req.URL.Path, time.Since(start))
	})
}
