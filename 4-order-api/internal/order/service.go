package order

import (
	"net/http"
	"orderApi/internal/product"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/logger"
)

type OrderRepositoryInterface interface {
	CreateOrder([]*product.Product, uint) (*Order, error)
	GetById(uint) (*Order, error)
	GetByUserId(uint) []Order
}

type OrderServiceDeps struct {
	OrderRepository OrderRepositoryInterface
}

type OrderService struct {
	OrderRepository OrderRepositoryInterface
}

func NewOrderService(deps *OrderServiceDeps) *OrderService {
	return &OrderService{
		OrderRepository: deps.OrderRepository,
	}
}

func (service *OrderService) CreateOrder(products []*product.Product, userId uint) (*Order, *apierrors.Error) {
	order, err := service.OrderRepository.CreateOrder(products, userId)

	if err != nil {
		logger.Error(err)
		return nil, apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return order, nil
}

func (service *OrderService) GetOrderById(orderId uint) (*Order, *apierrors.Error) {
	order, err := service.OrderRepository.GetById(uint(orderId))

	if err != nil {
		logger.Error(err)
		return nil, apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return order, nil
}

func (service *OrderService) GetOrdersByUserId(id uint) []Order {
	orders := service.OrderRepository.GetByUserId(id)

	return orders
}
