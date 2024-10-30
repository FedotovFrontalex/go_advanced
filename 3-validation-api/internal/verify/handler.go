package verify

import (
	"net/http"
	"os"
	"validationApi/configs"
	"validationApi/pkg/response"
)

type VerifyHandlerDeps struct {
		*configs.Config
}

type VerifyHandler struct {
		Deps VerifyHandlerDeps
}

func NewVerifyHandler(router *http.ServeMux, deps VerifyHandlerDeps) {
		handler := VerifyHandler{
				Deps: deps,
		}

		router.HandleFunc("POST /send", handler.send())
		router.HandleFunc("/verify/{hash}", handler.verify())
}

func (handler *VerifyHandler) send() http.HandlerFunc {
		return func (w http.ResponseWriter, req *http.Request) {
				payload := SendPayload{
						Email: os.Getenv("EMAIL"),
				}	
		
				response.Json(w, payload, 200)		
		}
}

func (handler *VerifyHandler) verify() http.HandlerFunc {
		return func (w http.ResponseWriter, req *http.Request) {
				payload := VerifyPayload{
						Email: os.Getenv("EMAIL"),
						Address: os.Getenv("ADDRESS"),
				}	

				response.Json(w, payload, 200)
		}
}
