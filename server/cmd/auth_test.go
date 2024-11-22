package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"server/internal/auth"
	"server/internal/user"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "test@test.ru",
		Password: "$2a$10$G0MfRZoTTUK8r.NVjUialOqVIPvveK24waYSMt7WJEGeEXKXAUdc6",
		Name:     "test  user",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email=?", "test@test.ru").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.ru",
		Password: "Byntuhfk",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

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

	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)

	if err != nil {
		t.Fatal(err)
	}

	token := resData.Token

	if token == "" {
		t.Fatal("Expected no empty string")
	}
}

func TestLoginFailed(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "1@1.ru",
		Password: "1",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 500 {
		t.Fatalf("expected %d got %d", 500, res.StatusCode)
	}
}
