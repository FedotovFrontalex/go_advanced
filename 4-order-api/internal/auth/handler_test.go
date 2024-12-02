package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"orderApi/configs"
	"orderApi/internal/auth"
	"orderApi/internal/user"
	apierrors "orderApi/pkg/apiErrors"
	"testing"
)

type MockAuthService struct{}

const USER_PHONE = "+79999999999"
const SESSION_ID = "i$I^imC0BtlVI*18$eSm@&"
const RightVerifyCode = 1111

func (s *MockAuthService) CreateSession(phone string) (string, error) {
	return SESSION_ID, nil
}

func (s *MockAuthService) VerifySession(sessionId string, code int) (*user.User, *apierrors.Error) {
	if code != RightVerifyCode {
		return nil, apierrors.NewError(http.StatusBadRequest, auth.ErrFailedCheckSession)
	}

	if sessionId != SESSION_ID {
		return nil, apierrors.NewError(http.StatusBadRequest, auth.ErrFailedCheckSession)
	}

	return &user.User{
		SessionId: SESSION_ID,
		Phone:     USER_PHONE,
	}, nil
}

func bootstrap() *auth.AuthHandler {
	return &auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "s8HDsT8z8e0Bn/lQwMhRpL+g685H2qVt0DdT/FxlYdA=",
			},
		},
		AuthService: &MockAuthService{},
	}
}

func TestAuthSuccess(t *testing.T) {
	handler := bootstrap()

	data, _ := json.Marshal(&auth.AuthRequest{
		Phone: USER_PHONE,
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth", reader)

	handler.Auth()(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d got %d", http.StatusCreated, w.Code)
	}
}

func TestAuthFailedWithWrongPhoneFormat(t *testing.T) {
	handler := bootstrap()

	data, _ := json.Marshal(&auth.AuthRequest{
		Phone: "1",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth", reader)

	handler.Auth()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestVerifyAuthSuccess(t *testing.T) {
	handler := bootstrap()

	data, _ := json.Marshal(&auth.AuthVerifyRequest{
		SessionId: SESSION_ID,
		Code:      RightVerifyCode,
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/verify", reader)

	handler.VerifyAuth()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code %d got %d", http.StatusOK, w.Code)
	}
}

func TestVerifyAuthFailedWithWrongSession(t *testing.T) {
	handler := bootstrap()

	data, _ := json.Marshal(&auth.AuthVerifyRequest{
		SessionId: "111",
		Code:      RightVerifyCode,
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/verify", reader)

	handler.VerifyAuth()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code %d got %d", http.StatusBadRequest, w.Code)
	}
}

func TestVerifyAuthFailedWithWrongCode(t *testing.T) {
	handler := bootstrap()

	data, _ := json.Marshal(&auth.AuthVerifyRequest{
		SessionId: SESSION_ID,
		Code:      5555,
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/verify", reader)

	handler.VerifyAuth()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected status code %d got %d", http.StatusBadRequest, w.Code)
	}
}
