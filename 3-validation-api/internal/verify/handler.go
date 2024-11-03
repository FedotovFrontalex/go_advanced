package verify

import (
	"net/http"
	"os"
	"strings"
	"validationApi/configs"
	"validationApi/internal/email"
	"validationApi/pkg/logger"
	"validationApi/pkg/request"
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
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("SEND")

		body, err := request.HandleBody[SendRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		err = email.SendVerifyEmail(strings.Fields(body.Email), handler.Deps.Config)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 500)
			return
		}

		payload := SendPayload{
			Email: os.Getenv("EMAIL"),
		}

		response.Json(w, payload, 200)
	}
}

func (handler *VerifyHandler) verify() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hash := req.PathValue("hash")

		status, err := email.CheckEmail(hash)

		if err != nil {
			response.Json(w, err.Error(), status)
			return
		}

		response.Json(w, "email verified", status)
	}
}
