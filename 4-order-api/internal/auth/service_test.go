package auth_test

import (
	"net/http"
	"orderApi/configs"
	"orderApi/internal/auth"
	"orderApi/internal/user"
	apierrors "orderApi/pkg/apiErrors"
	"testing"
)

type MockUserService struct{}

func (s *MockUserService) FindUserByPhone(phone string) (*user.User, *apierrors.Error) {

	if phone != USER_PHONE {
		return nil, apierrors.NewError(404, "user not found")
	}

	return &user.User{
		Phone: USER_PHONE,
	}, nil
}

func (s *MockUserService) UpdateUser(userData *user.User) (*user.User, *apierrors.Error) {
	return &user.User{
		Phone:     userData.Phone,
		SessionId: userData.SessionId,
	}, nil
}

func (s *MockUserService) CreateUser(phone string) (*user.User, *apierrors.Error) {
	return &user.User{
		Phone:     phone,
		SessionId: SESSION_ID,
	}, nil
}

func (s *MockUserService) FindBySessionId(sessionId string) (*user.User, *apierrors.Error) {
	if sessionId != SESSION_ID {
		return nil, apierrors.NewError(404, "session not found")
	}

	return &user.User{
		Phone:     USER_PHONE,
		SessionId: sessionId,
	}, nil
}

func TestServiceCreateSessionByExistdUserSuccess(t *testing.T) {
	const sessionIdLength = 22

	conf := &configs.Config{
		SessionIdLength: sessionIdLength,
	}

	service := auth.NewAuthService(&auth.AuthServiceDeps{
		Config:      conf,
		UserService: &MockUserService{},
	})

	session, err := service.CreateSession(USER_PHONE)

	if err != nil {
		t.Fatal(err)
		return
	}

	if len(session) != sessionIdLength {
		t.Fatalf("SessionId length expected %d  got %d", sessionIdLength, len(session))
	}
}

func TestServiceCreateSessionByNewUserSuccess(t *testing.T) {
	const sessionIdLength = 22

	conf := &configs.Config{
		SessionIdLength: sessionIdLength,
	}

	service := auth.NewAuthService(&auth.AuthServiceDeps{
		Config:      conf,
		UserService: &MockUserService{},
	})

	session, err := service.CreateSession("81111111111")

	if err != nil {
		t.Fatal(err)
		return
	}

	if len(session) != sessionIdLength {
		t.Fatalf("SessionId length expected %d  got %d", sessionIdLength, len(session))
	}
}

func TestServiceVerifySessionByExistedUserSuccess(t *testing.T) {
	const sessionIdLength = 22

	conf := &configs.Config{
		SessionIdLength: sessionIdLength,
	}

	service := auth.NewAuthService(&auth.AuthServiceDeps{
		Config:      conf,
		UserService: &MockUserService{},
	})

	user, err := service.VerifySession(SESSION_ID, RightVerifyCode)

	if err != nil {
		t.Fatal(err)
		return
	}

	if user.SessionId != SESSION_ID {
		t.Fatal("session ids not matched")
		return
	}
}

func TestServiceVerifySessionByExistedUserFailed(t *testing.T) {
	const sessionIdLength = 22

	conf := &configs.Config{
		SessionIdLength: sessionIdLength,
	}

	service := auth.NewAuthService(&auth.AuthServiceDeps{
		Config:      conf,
		UserService: &MockUserService{},
	})

	user, err := service.VerifySession(SESSION_ID, 2222)

	if err == nil {
		t.Fatal(err)
		return
	}

	if err.GetStatus() != http.StatusBadRequest {
		t.Fatalf("Expected status code %d got %d", http.StatusBadRequest, err.GetStatus())
	}

	if err.Error() != auth.ErrFailedCheckSession {
		t.Fatalf("Expected error message %s got %s", auth.ErrFailedCheckSession, err.Error())
	}

	if user != nil {
		t.Fatal("user must be nil")
	}
}

func TestServiceVerifySessionByWrongSessionIdUserFailed(t *testing.T) {
	const sessionIdLength = 22

	conf := &configs.Config{
		SessionIdLength: sessionIdLength,
	}

	service := auth.NewAuthService(&auth.AuthServiceDeps{
		Config:      conf,
		UserService: &MockUserService{},
	})

	user, err := service.VerifySession("wrongSession", RightVerifyCode)

	if err == nil {
		t.Fatal(err)
		return
	}

	if err.GetStatus() != http.StatusBadRequest {
		t.Fatalf("Expected status code %d got %d", http.StatusBadRequest, err.GetStatus())
	}

	if err.Error() != auth.ErrFailedCheckSession {
		t.Fatalf("Expected error message %s got %s", auth.ErrFailedCheckSession, err.Error())
	}

	if user != nil {
		t.Fatal("user must be nil")
	}
}
