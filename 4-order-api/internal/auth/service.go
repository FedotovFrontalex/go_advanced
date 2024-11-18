package auth

import (
	"errors"
	"orderApi/configs"
	"orderApi/internal/user"
	"orderApi/pkg/logger"
)

type AuthServiceDeps struct {
	*configs.Config
	UserRepository *user.UserRepository
}

type AuthService struct {
	UserRepository  *user.UserRepository
	SessionIdLength int
}

func NewAuthService(deps *AuthServiceDeps) *AuthService {
	return &AuthService{
		UserRepository:  deps.UserRepository,
		SessionIdLength: deps.Config.SessionIdLength,
	}
}

func (service *AuthService) CreateSession(phone string) (string, error) {
	existedUser, _ := service.UserRepository.FindByPhone(phone)

	var userData *user.User

	if existedUser == nil {
		userData = user.NewUser(phone)
	} else {
		userData = existedUser
	}

	for {
		userData.CreateSessionId(service.SessionIdLength)
		userWithSameSessionId, _ := service.UserRepository.FindBySessionId(userData.SessionId)

		if userWithSameSessionId == nil {
			break
		}
	}

	if existedUser == nil {
		user, err := service.UserRepository.Create(userData)

		if err != nil {
			return "", err
		}

		return user.SessionId, nil
	} else {
		user, err := service.UserRepository.Update(userData)

		if err != nil {
			return "", nil
		}

		return user.SessionId, nil
	}
}

func (service *AuthService) VerifySession(sessionId string, code int) (*user.User, error) {
	user, err := service.UserRepository.FindBySessionId(sessionId)

	if err != nil {
		return nil, errors.New(ErrFailedCheckSession)
	}

	if code != 1111 {
		user.SessionId = " "
		_, err = service.UserRepository.Update(user)

		if err != nil {
			logger.Error(errors.New(ErrFailedClearSession))
		}
		return nil, errors.New(ErrFailedCheckSession)
	}

	return user, nil
}
