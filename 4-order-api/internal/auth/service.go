package auth

import (
	"errors"
	"net/http"
	"orderApi/configs"
	"orderApi/internal/user"
	apierrors "orderApi/pkg/apiErrors"
	"orderApi/pkg/logger"
)

type UserServiceInterface interface {
	FindUserByPhone(string) (*user.User, *apierrors.Error)
	CreateUser(string) (*user.User, *apierrors.Error)
	UpdateUser(*user.User) (*user.User, *apierrors.Error)
	FindBySessionId(string) (*user.User, *apierrors.Error)
}

type AuthServiceDeps struct {
	*configs.Config
	UserService UserServiceInterface
}

type AuthService struct {
	UserService     UserServiceInterface
	SessionIdLength int
}

func NewAuthService(deps *AuthServiceDeps) *AuthService {
	return &AuthService{
		UserService:     deps.UserService,
		SessionIdLength: deps.Config.SessionIdLength,
	}
}

func (service *AuthService) CreateSession(phone string) (string, error) {
	existedUser, _ := service.UserService.FindUserByPhone(phone)

	if existedUser == nil {
		user, apierr := service.UserService.CreateUser(phone)

		if apierr != nil {
			return "", apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		return user.SessionId, nil
	}

	existedUser.CreateSessionId(service.SessionIdLength)
	user, apierr := service.UserService.UpdateUser(existedUser)

	if apierr != nil {
		return "", apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return user.SessionId, nil
}

func (service *AuthService) VerifySession(sessionId string, code int) (*user.User, *apierrors.Error) {
	user, err := service.UserService.FindBySessionId(sessionId)

	if err != nil {
		return nil, apierrors.NewError(http.StatusBadRequest, ErrFailedCheckSession)
	}

	if code != 1111 {
		user.SessionId = " "
		_, err = service.UserService.UpdateUser(user)

		if err != nil {
			logger.Error(errors.New(ErrFailedClearSession))
		}
		return nil, apierrors.NewError(http.StatusBadRequest, ErrFailedCheckSession)
	}

	return user, nil
}
