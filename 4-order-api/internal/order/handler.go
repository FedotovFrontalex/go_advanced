package order

import (
	"net/http"
	"orderApi/configs"
	"orderApi/internal/product"
	"orderApi/internal/user"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/middleware"
	"orderApi/pkg/request"
	"orderApi/pkg/response"
	"strconv"
)

type OrderServiceInterface interface {
	CreateOrder([]*product.Product, uint) (*Order, *apierrors.Error)
	GetOrderById(uint) (*Order, *apierrors.Error)
	GetOrdersByUserId(uint) []Order
}

type ProductServiceInterface interface {
	GetProductsSliceById([]string) []*product.Product
}

type UserServiceInterface interface {
	FindBySessionId(string) (*user.User, *apierrors.Error)
}

type OrderHandlerDeps struct {
	Config         *configs.Config
	OrderService   OrderServiceInterface
	ProductService ProductServiceInterface
	UserService    UserServiceInterface
}

type OrderHandler struct {
	OrderService   OrderServiceInterface
	ProductService ProductServiceInterface
	UserService    UserServiceInterface
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderService:   deps.OrderService,
		ProductService: deps.ProductService,
		UserService:    deps.UserService,
	}

	router.Handle("POST /order", middleware.IsAuthed(handler.CreateOrder(), deps.Config))
	router.Handle("/order/{id}", middleware.IsAuthed(handler.GetById(), deps.Config))
	router.Handle("/my-orders", middleware.IsAuthed(handler.GetByUserId(), deps.Config))
}

func (handler *OrderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionId, ok := req.Context().Value(middleware.ContextSessionIdKey).(string)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, apierr := handler.UserService.FindBySessionId(sessionId)

		if apierr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		body, err := request.HandleBody[OrderRequest](&w, req)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		products := handler.ProductService.GetProductsSliceById(body.Products)

		order, apiErr := handler.OrderService.CreateOrder(products, user.ID)

		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.GetStatus())
			return
		}

		response.Json(w, order, 201)
	}
}

func (handler *OrderHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionId, ok := req.Context().Value(middleware.ContextSessionIdKey).(string)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		orderId, err := strconv.ParseUint(req.PathValue("id"), 10, 64)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user, apierr := handler.UserService.FindBySessionId(sessionId)

		if apierr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		order, errApi := handler.OrderService.GetOrderById(uint(orderId))

		if errApi != nil {
			http.Error(w, errApi.Error(), errApi.GetStatus())
		}

		if order.UserId != user.ID {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		response.Json(w, order, 200)
	}
}

func (handler *OrderHandler) GetByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		sessionId, ok := req.Context().Value(middleware.ContextSessionIdKey).(string)

		if !ok {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		user, err := handler.UserService.FindBySessionId(sessionId)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		orders := handler.OrderService.GetOrdersByUserId(user.ID)

		response.Json(w, orders, 200)
	}
}
