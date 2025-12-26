package ports

import "apitest/internal/core/user"

type UserUseCase interface {
	CreateUser(user *user.AppUser) error
	GetUsersByIds(ids ...int) ([]*user.AppUser, error)
	GetUserById(id int) (*user.AppUser, error)
	GetUsersByPage(cursor int, limit int) ([]*user.AppUser, error)
}
