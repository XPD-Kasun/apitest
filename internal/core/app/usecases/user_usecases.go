package usecases

import (
	"apitest/internal/core/user"
)

func NewUserUseCaseImpl(uc user.UserService) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		userSvc: uc,
	}
}

type UserUseCaseImpl struct {
	userSvc user.UserService
}

// GetUsersByPage implements [ports.UserUseCase].
func (u *UserUseCaseImpl) GetUsersByPage(cursor int, limit int) ([]*user.AppUser, error) {
	return u.userSvc.GetByPage(cursor, limit)
}

// GetUserById implements [ports.UserUseCase].
func (u *UserUseCaseImpl) GetUserById(id int) (*user.AppUser, error) {
	return u.userSvc.GetById(id)
}

// GetUsersByIds implements [ports.UserUseCase].
func (u *UserUseCaseImpl) GetUsersByIds(ids ...int) ([]*user.AppUser, error) {
	return u.userSvc.GetByIds(ids...)
}

// CreateUser implements ports.UserUseCase.
func (u *UserUseCaseImpl) CreateUser(user *user.AppUser) error {
	return u.userSvc.CreateNew(user)
}

// vs implmenter
// var s ports.UserUseCase = &UserUseCaseImpl{}
