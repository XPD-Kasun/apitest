package dataloaders

import (
	"apitest/internal/core/app/ports"
	"apitest/internal/core/user"
	"context"
)

type UserDataloader struct {
	userUC ports.UserUseCase
}

func cloneError(err error, count int) []error {
	errorset := make([]error, count)
	for i := 0; i < count; i++ {
		errorset[i] = err
	}
	return errorset
}

func (ud *UserDataloader) GetUsers(ctx context.Context, ids []int) ([]*user.AppUser, []error) {
	if ctx.Err() != nil {
		return nil, cloneError(ctx.Err(), len(ids))
	}
	users, err := ud.userUC.GetUsersByIds(ids...)
	errorset := cloneError(err, len(ids))
	return users, errorset
}

func (ud *UserDataloader) GetUser(ctx context.Context, id int) (*user.AppUser, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return ud.userUC.GetUserById(id)
}
