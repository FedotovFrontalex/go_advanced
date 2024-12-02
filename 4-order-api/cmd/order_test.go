package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"orderApi/internal/auth"
	"orderApi/internal/order"
	"orderApi/internal/product"
	"testing"

	"gorm.io/gorm"
)

func initProductData(db *gorm.DB) {
	db.Create(&product.Product{
		Name:        "product 1",
		Description: "Desc 1",
		Model: gorm.Model{
			ID: 1,
		},
	})

	db.Create(&product.Product{
		Name:        "product 2",
		Description: "Desc 2",
		Model: gorm.Model{
			ID: 2,
		},
	})

	db.Create(&product.Product{
		Name:        "product 3",
		Description: "Desc 3",
		Model: gorm.Model{
			ID: 3,
		},
	})
}

func removeProductData(db *gorm.DB) {
	var products = []product.Product{
		{Model: gorm.Model{ID: 1}},
		{Model: gorm.Model{ID: 2}},
		{Model: gorm.Model{ID: 3}},
	}

	db.Unscoped().
		Delete(&products)
}

func removeOrderData(db *gorm.DB, id uint) {
	db.Unscoped().
		Table("order_products").
		Where("order_id=?", id).
		Delete(&product.Product{})

	db.Unscoped().
		Where("id=?", id).
		Delete(&order.Order{})
}

func TestCreateOrderSuccess(t *testing.T) {
	db := initDb()

	initUserData(db)
	defer removeUserData(db)

	initProductData(db)
	defer removeProductData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.AuthRequest{
		Phone: "+79999999999",
	})

	res, err := http.Post(ts.URL+"/auth", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	var resData auth.AuthResponse
	err = json.Unmarshal(body, &resData)

	if err != nil {
		t.Fatal(err)
	}

	sessionId := resData.SessionId

	if sessionId == "" {
		t.Fatal("Expected no Empty string")
	}

	if sessionId == "SESSION_ID" {
		t.Fatal("Must be created new session id")
	}

	dataVerify, _ := json.Marshal(&auth.AuthVerifyRequest{
		SessionId: sessionId,
		Code:      1111,
	})

	res, err = http.Post(ts.URL+"/auth/verify", "application/json", bytes.NewReader(dataVerify))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	var resVerifyData auth.AuthVerificationResponse
	err = json.Unmarshal(body, &resVerifyData)

	if err != nil {
		t.Fatal(err)
	}

	token := resVerifyData.Token

	if token == "" {
		t.Fatal("Token must be non empty string")
	}

	orderData, _ := json.Marshal(&order.OrderRequest{
		Products: []string{"1", "3"},
	})

	req, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(orderData))

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-type", "Application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status code expected %d got %d", http.StatusCreated, res.StatusCode)
	}

	body, err = io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	var order *order.Order
	err = json.Unmarshal(body, &order)

	if err != nil {
		t.Fatal(err)
	}

	defer removeOrderData(db, order.ID)

	if len(order.Products) != 2 {
		t.Fatalf("Expected count products %d got %d", 2, len(order.Products))
	}
}
