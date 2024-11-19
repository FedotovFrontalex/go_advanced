package product

import (
	"net/http"
	"orderApi/configs"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/logger"
	"orderApi/pkg/middleware"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
	"strconv"

	"github.com/lib/pq"
)

type ProductServiceInterface interface {
	CreateProduct(string, string, pq.StringArray) (*Product, *apierrors.Error)
	UpdateProduct(uint, string, string, pq.StringArray) (*Product, *apierrors.Error)
	DeleteProduct(uint) *apierrors.Error
	GetProductById(uint) (*Product, *apierrors.Error)
}

type ProductHandlerDeps struct {
	*configs.Config
	ProductService ProductServiceInterface
}

type ProductHandler struct {
	*configs.Config
	ProductService ProductServiceInterface
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	logger.Message("initialize routes: product")
	handler := &ProductHandler{
		Config:         deps.Config,
		ProductService: deps.ProductService,
	}

	router.Handle("POST /product", middleware.IsAuthed(handler.CreateProduct(), deps.Config))
	router.Handle("PATCH /product/{id}", middleware.IsAuthed(handler.UpdateProduct(), deps.Config))
	router.Handle("DELETE /product/{id}", middleware.IsAuthed(handler.DeleteProduct(), deps.Config))
	router.HandleFunc("/product/{id}", handler.GetProductById())
}

func (handler *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("Add product")

		_, ok := req.Context().Value(middleware.ContextSessionIdKey).(string)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		body, err := request.HandleBody[ProductCreateRequest](&w, req)

		if err != nil {
			logger.Error(err)
			response.Json(w, err.Error(), 402)
			return
		}

		product, apierr := handler.ProductService.CreateProduct(body.Name, body.Description, body.Images)

		if apierr != nil {
			http.Error(w, apierr.Error(), apierr.GetStatus())
			return
		}

		response.Json(w, product, 201)
	}
}

func (handler *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("UpdateProduct")

		_, ok := req.Context().Value(middleware.ContextSessionIdKey).(string)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		id, err := strconv.ParseUint(req.PathValue("id"), 10, 32)

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

		product, apierr := handler.ProductService.UpdateProduct(uint(id), body.Name, body.Description, body.Images)

		if apierr != nil {
			http.Error(w, apierr.Error(), apierr.GetStatus())
			return
		}

		response.Json(w, product, 201)
	}
}

func (handler *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		logger.Message("DeleteProduct")

		_, ok := req.Context().Value(middleware.ContextSessionIdKey).(string)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		id, err := strconv.ParseUint(req.PathValue("id"), 10, 32)

		if err != nil {
			logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		apierr := handler.ProductService.DeleteProduct(uint(id))

		if apierr != nil {
			http.Error(w, apierr.Error(), apierr.GetStatus())
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

		product, apierr := handler.ProductService.GetProductById(uint(id))

		if apierr != nil {
			http.Error(w, apierr.Error(), apierr.GetStatus())
			return
		}

		response.Json(w, product, 200)
	}
}
