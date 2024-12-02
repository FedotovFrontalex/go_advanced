package order_test

import (
	"orderApi/internal/order"
	"orderApi/internal/product"
	"testing"
)

type MockOrderRepository struct{}

func (r *MockOrderRepository) CreateOrder(products []*product.Product, userId uint) (*order.Order, error) {
	return &order.Order{
		UserId:   userId,
		Products: products,
	}, nil
}

func (r *MockOrderRepository) GetById(orderId uint) (*order.Order, error) {
	product1 := &product.Product{
		Name:        "Product1",
		Description: "Desc 1",
	}

	product2 := &product.Product{
		Name:        "Product2",
		Description: "Desc 2",
	}

	products := []*product.Product{product1, product2}

	return &order.Order{
		Products: products,
	}, nil
}

func (r *MockOrderRepository) GetByUserId(userId uint) []order.Order {
	product1 := &product.Product{
		Name:        "Product1",
		Description: "Desc 1",
	}

	product2 := &product.Product{
		Name:        "Product2",
		Description: "Desc 2",
	}

	products := []*product.Product{product1, product2}

	order1 := order.Order{
		UserId:   userId,
		Products: products,
	}

	return []order.Order{order1}
}

func TestCreateOrderSuccess(t *testing.T) {
	service := &order.OrderService{
		OrderRepository: &MockOrderRepository{},
	}

	product1 := &product.Product{
		Name:        "Product1",
		Description: "Desc 1",
	}

	product2 := &product.Product{
		Name:        "Product2",
		Description: "Desc 2",
	}

	products := []*product.Product{product1, product2}

	order, err := service.CreateOrder(products, USER_ID)

	if err != nil {
		t.Fatal(err)
		return
	}

	if len(order.Products) != len(products) {
		t.Fatalf("Products lenght expected %d got %d", len(products), len(order.Products))
	}
}

func TestGetByUserIdSuccess(t *testing.T) {
	service := &order.OrderService{
		OrderRepository: &MockOrderRepository{},
	}

	orders := service.GetOrdersByUserId(USER_ID)

	for _, val := range orders {
		if val.UserId != USER_ID {
			t.Fatalf("User id expected %d got %d", USER_ID, val.UserId)

		}
	}
}
