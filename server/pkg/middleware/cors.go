package middleware

import "net/http"

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")

		if origin == "" {
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
