package sqldto

import (
	"apitest/internal/core/user"

	"github.com/uptrace/bun"
)

type UserSqlDto struct {
	bun.BaseModel `bun:"table:appuser"`

	Id        int
	Username  string
	Firstname string
	Lastname  string
	Email     string
}

func (u *UserSqlDto) ToCoreUser() *user.AppUser {
	return &user.AppUser{
		Id:        u.Id,
		UserName:  u.Username,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
	}
}
