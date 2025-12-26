package user

import "apitest/internal/core/common/baserepo"

type AppUserRepo interface {
	baserepo.MultiIdGetter[*AppUser, int]
	baserepo.SingleIdGetter[*AppUser, int]
	baserepo.Inserter[*AppUser]
	baserepo.PaginatedGetter[*AppUser, int]
	GetPasswordHash(id int) (string, error)
	GetByUserName(username string) (*AppUser, error)
}
