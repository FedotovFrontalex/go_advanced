package order_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"orderApi/internal/order"
	"orderApi/internal/product"
	"orderApi/internal/user"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/middleware"
	"strconv"
	"testing"

	"gorm.io/gorm"
)

const EXISTED_SESSION = "existedSessionId"
const NON_EXSITED_SESSION = "nonExistedSession"
const ORDER_ID = 1
const USER_ID = 1
const ORDER_ID_2 = 2
const USER_ID_2 = 2

type MockOrderService struct{}
type MockProductService struct{}
type MockUserService struct{}

func (service *MockOrderService) CreateOrder(products []*product.Product, userId uint) (*order.Order, *apierrors.Error) {
	return &order.Order{
		Products: products,
		UserId:   userId,
	}, nil
}

func (service *MockOrderService) GetOrderById(orderId uint) (*order.Order, *apierrors.Error) {
	if orderId != ORDER_ID && orderId != ORDER_ID_2 {
		return nil, apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	orderData := &order.Order{
		Model: gorm.Model{
			ID: orderId,
		},
		UserId: USER_ID,
	}

	if orderId == ORDER_ID_2 {
		orderData.UserId = USER_ID_2
	}

	return orderData, nil
}

func (service *MockOrderService) GetOrdersByUserId(userId uint) []order.Order {
	order1 := order.Order{
		UserId: userId,
	}

	order2 := order.Order{
		UserId: userId,
	}

	return []order.Order{order1, order2}
}

func (service *MockProductService) GetProductsSliceById(productsIds []string) []*product.Product {
	products := []*product.Product{}

	for _, val := range productsIds {
		id, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			products = append(products, &product.Product{
				Model: gorm.Model{ID: uint(id)},
			})
		}
	}

	return products
}

func (service *MockUserService) FindBySessionId(sessionId string) (*user.User, *apierrors.Error) {
	if sessionId == NON_EXSITED_SESSION {
		return nil, apierrors.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return &user.User{
		SessionId: sessionId,
		Model:     gorm.Model{ID: USER_ID},
	}, nil
}

func bootstrap() *order.OrderHandler {
	orderService := &MockOrderService{}
	productService := &MockProductService{}
	userService := &MockUserService{}

	return &order.OrderHandler{
		OrderService:   orderService,
		ProductService: productService,
		UserService:    userService,
	}
}

func TestHandlerCreateOrderSuccess(t *testing.T) {
	handler := bootstrap()

	productsIds := []string{"1", "2"}

	data, _ := json.Marshal(&order.OrderRequest{
		Products: productsIds,
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/order", reader)

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, EXISTED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.CreateOrder()(w, reqWithCtx)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d got %d", http.StatusCreated, w.Code)
	}

	body, err := io.ReadAll(w.Result().Body)

	if err != nil {
		t.Fatal(err)
	}

	var orderData order.Order
	json.Unmarshal(body, &orderData)

	if len(productsIds) != len(orderData.Products) {
		t.Fatalf("expected products count in order %d got %d", len(productsIds), len(orderData.Products))
	}

	if orderData.UserId != USER_ID {
		t.Fatalf("Expected user id %d got %d", USER_ID, orderData.UserId)
	}
}

func TestHandlerCreateOrderFailedNoSession(t *testing.T) {
	handler := bootstrap()

	productsIds := []string{"1", "2"}

	data, _ := json.Marshal(&order.OrderRequest{
		Products: productsIds,
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/order", reader)

	handler.CreateOrder()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandlerCreateOrderFailedWrongSession(t *testing.T) {
	handler := bootstrap()

	productsIds := []string{"1", "2"}

	data, _ := json.Marshal(&order.OrderRequest{
		Products: productsIds,
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/order", reader)

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, NON_EXSITED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.CreateOrder()(w, reqWithCtx)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status %d got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandlerGetByIdSuccess(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	req.SetPathValue("id", strconv.Itoa(ORDER_ID))

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, EXISTED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.GetById()(w, reqWithCtx)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d got %d", http.StatusOK, w.Code)
	}

	body, err := io.ReadAll(w.Result().Body)

	if err != nil {
		t.Fatal(err)
	}

	var orderData order.Order
	json.Unmarshal(body, &orderData)

	if orderData.ID != ORDER_ID {
		t.Fatalf("expected order id %d got %d", ORDER_ID, orderData.ID)
	}

	if orderData.UserId != USER_ID {
		t.Fatalf("expected user id %d got %d", USER_ID, orderData.UserId)
	}
}

func TestHandlerGetByIdFailedByAnotherUserOrder(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	req.SetPathValue("id", strconv.Itoa(ORDER_ID_2))

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, EXISTED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.GetById()(w, reqWithCtx)

	if w.Code != http.StatusForbidden {
		t.Fatalf("Expected status %d got %d", http.StatusForbidden, w.Code)
	}
}

func TestHandlerGetByIdFailedUnautorized(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	req.SetPathValue("id", strconv.Itoa(ORDER_ID_2))

	handler.GetById()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandlerGetByIdFailedWrongSession(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/order", nil)
	req.SetPathValue("id", strconv.Itoa(ORDER_ID_2))

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, NON_EXSITED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.GetById()(w, reqWithCtx)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status %d got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandlerGetByUserIdSuccess(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/my-orders", nil)

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, EXISTED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.GetByUserId()(w, reqWithCtx)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d got %d", http.StatusOK, w.Code)
	}

	body, err := io.ReadAll(w.Result().Body)

	if err != nil {
		t.Fatal(err)
	}

	var orders []order.Order
	json.Unmarshal(body, &orders)

	if orders == nil {
		t.Fatal("Expected list of orders")
	}
}

func TestHandlerGetByUserIdFailedNoSession(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/my-orders", nil)

	handler.GetByUserId()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status %d got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandlerGetByUserIdWrongSession(t *testing.T) {
	handler := bootstrap()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/my-orders", nil)

	ctx := context.WithValue(req.Context(), middleware.ContextSessionIdKey, NON_EXSITED_SESSION)
	reqWithCtx := req.WithContext(ctx)

	handler.GetByUserId()(w, reqWithCtx)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected status %d got %d", http.StatusInternalServerError, w.Code)
	}
}
