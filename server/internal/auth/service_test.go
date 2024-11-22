package auth_test

import (
	"server/internal/auth"
	"server/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (mock *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "a@a.ru",
	}, nil
}

func (mock *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterServiceSuccess(t *testing.T) {
	const emailReg = "testReg@test.ru"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register(emailReg, "1", "Alex")

	if err != nil {
		t.Fatal(err)
	}

	if email != emailReg {
		t.Fatalf("Expected %s got %s", emailReg, email)
	}
}
