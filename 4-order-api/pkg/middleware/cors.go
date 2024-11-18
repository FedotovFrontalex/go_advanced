package middleware

import (
	"net/http"
	"orderApi/configs"
	"strings"
)

func InitCors(securityConfig configs.SecurityConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		domains := strings.Split(securityConfig.Domains, ",")
		domainsMap := make(map[string]bool)

		for _, value := range domains {
			domainsMap[value] = true
		}

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")

			_, ok := domainsMap[origin]

			if !ok {
				next.ServeHTTP(w, req)
				return
			}

			header := w.Header()
			header.Set("Access-Control-Allow-Origin", origin)
			header.Set("Access-Control-Allow-Credentials", "true")

			if req.Method == http.MethodOptions {
				header.Set("Access-Control-Allow-Metods", "GET,PUT,POST,DELETE,HEAD,PATCH")
				header.Set("Access-Control-Allow0Headers", "autorization,content-type,content-length")
				header.Set("Access-Control-Max-Age", "86400")
			}

			next.ServeHTTP(w, req)
		})
	}
}
