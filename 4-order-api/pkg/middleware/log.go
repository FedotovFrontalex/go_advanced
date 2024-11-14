package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, req)
		log.SetFormatter(&log.JSONFormatter{})
		log.WithFields(log.Fields{
			"statusCode": wrapper.StatusCode,
			"method":     req.Method,
			"path":       req.URL.Path,
			"duration":   time.Since(start),
		}).Info("api call")
	})
}
