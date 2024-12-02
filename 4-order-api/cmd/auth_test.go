package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"orderApi/internal/auth"
	"orderApi/internal/user"
	"os"
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

func initUserData(db *gorm.DB) {
	db.Create(&user.User{
		Phone:     "+79999999999",
		SessionId: "SESSION_ID",
	})
}

func removeUserData(db *gorm.DB) {
	db.Unscoped().
		Where("phone=?", "+79999999999").
		Delete(&user.User{})
}

func TestLoginNewUserSuccess(t *testing.T) {
	// db := initDb()

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

	var resVerifyData auth.AuthVerificationResponse
	err = json.Unmarshal(body, &resVerifyData)

	if err != nil {
		t.Fatal(err)
	}

	token := resVerifyData.Token

	if token == "" {
		t.Fatal("Token must be non empty string")
	}
}

func TestLoginExistedUserSuccess(t *testing.T) {
	db := initDb()
	initUserData(db)
	defer removeUserData(db)

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
}
