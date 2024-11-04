package product

import (
	"net/http"
	"orderApi/configs"
	"orderApi/pkg/logger"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
)

type ProductHandlerDeps struct {
	*configs.Config
}

type ProductHandler struct {
	*configs.Config
}

func NewProductHandler(router *http.ServeMux, conf *ProductHandlerDeps) {
	logger.Message("initialize routes: product")
	handler := ProductHandler{
		Config: conf.Config,
	}

	router.HandleFunc("POST /product/addProduct", handler.addProduct())
}

func (handler *ProductHandler) addProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("Add product")
		data, err := request.HandleBody[AddProductRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		response.Json(w, data, 201)
	}
}
