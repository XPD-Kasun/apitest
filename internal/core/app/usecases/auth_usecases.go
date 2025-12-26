package usecases

import (
	"apitest/internal/core/app/dto"
	"apitest/internal/core/user"
)

type AuthUseCaseImpl struct {
	userSvc user.UserService
}

func NewAuthUseCaseImpl(userSvc user.UserService) *AuthUseCaseImpl {

	return &AuthUseCaseImpl{
		userSvc: userSvc,
	}

}

func (a *AuthUseCaseImpl) LoginFromPassword(username, password string) (string, error) {
	return a.userSvc.Login(username, password)
}

func (a *AuthUseCaseImpl) CreateUser(user *dto.UserDTO) error {
	return a.userSvc.CreateNew(user.ToAppUser())
}
