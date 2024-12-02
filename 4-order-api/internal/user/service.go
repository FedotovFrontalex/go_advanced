package user

import (
	"net/http"
	"orderApi/configs"
	apierrors "orderApi/pkg/apiErrors"
)

type UserRepositoryInterface interface {
	Create(*User) (*User, error)
	Update(*User) (*User, error)
	FindByPhone(string) (*User, error)
	FindBySessionId(string) (*User, error)
}

type UserServiceDeps struct {
	*configs.Config
	UserRepository UserRepositoryInterface
}

type UserService struct {
	UserRepository  UserRepositoryInterface
	SessionIdLength int
}

func NewUserService(deps *UserServiceDeps) *UserService {
	return &UserService{
		UserRepository:  deps.UserRepository,
		SessionIdLength: deps.Config.SessionIdLength,
	}
}

func (service *UserService) CreateUser(phone string) (*User, *apierrors.Error) {
	user := NewUser(phone)
	user.CreateSessionId(service.SessionIdLength)

	createdUser, err := service.UserRepository.Create(user)

	if err != nil {
		return nil, apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return createdUser, nil
}

func (service *UserService) UpdateUser(user *User) (*User, *apierrors.Error) {
	updatedUser, err := service.UserRepository.Update(user)

	if err != nil {
		return nil, apierrors.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	return updatedUser, nil
}

func (service *UserService) FindUserByPhone(phone string) (*User, *apierrors.Error) {
	existedUser, err := service.UserRepository.FindByPhone(phone)

	if err != nil {
		return nil, apierrors.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return existedUser, nil
}

func (service *UserService) FindBySessionId(sessionId string) (*User, *apierrors.Error) {
	user, err := service.UserRepository.FindBySessionId(sessionId)

	if err != nil {
		return nil, apierrors.NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return user, nil
}
