package main

import (
	"math/rand/v2"
	"net/http"
	"strconv"
)

type ApiHandler struct{}

func NewApiHandler(router *http.ServeMux) {
		handler := &ApiHandler{}
		router.HandleFunc("/random", handler.randomByte())
}

func (handler *ApiHandler) randomByte() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
				val := rand.N(6) + 1
				w.Write([]byte(strconv.Itoa(val)))
		}
}
