package product

import (
	"net/http"
	"orderApi/configs"
	"orderApi/pkg/logger"
	"orderApi/pkg/middleware"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	*configs.Config
	ProductRepository *ProductRepository
}

type ProductHandler struct {
	*configs.Config
	ProductRepository *ProductRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	logger.Message("initialize routes: product")
	handler := &ProductHandler{
		Config:            deps.Config,
		ProductRepository: deps.ProductRepository,
	}

	router.Handle("POST /product", middleware.IsAuthed(handler.CreateProduct(), deps.Config))
	router.Handle("PATCH /product/{id}", middleware.IsAuthed(handler.UpdateProduct(), deps.Config))
	router.HandleFunc("/product/{id}", handler.GetProductById())
	router.HandleFunc("DELETE /product/{id}", handler.DeleteProduct())
}

func (handler *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("Add product")
		body, err := request.HandleBody[ProductCreateRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		product := NewProduct(body.Name, body.Description, body.Images)
		result, err := handler.ProductRepository.Create(product)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, result, 201)
	}
}

func (handler *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("UpdateProduct")

		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		body, err := request.HandleBody[ProductUpdateRequest](&w, req)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.ProductRepository.GetProductById(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := handler.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Images:      body.Images,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, result, 201)
	}
}

func (handler *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("DeleteProduct")

		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.ProductRepository.GetProductById(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.ProductRepository.Delete(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, "delete success", 200)
	}
}

func (handler *ProductHandler) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("GetProductById")

		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := handler.ProductRepository.GetProductById(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response.Json(w, result, 200)
	}
}
