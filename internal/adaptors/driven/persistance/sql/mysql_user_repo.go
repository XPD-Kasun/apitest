package sql

import (
	"apitest/internal/adaptors/driven/persistance/sql/sqldto"
	"apitest/internal/core/common/baserepo"
	"apitest/internal/core/common/filters"
	"apitest/internal/core/common/funcs"
	"apitest/internal/core/user"
	"apitest/internal/logger"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
)

func NewMySqlUserRepo(db *bun.DB) (*MySqlUserRepo, error) {
	return &MySqlUserRepo{db: db}, nil
}

type MySqlUserRepo struct {
	db *bun.DB
}

// GetByIds implements [user.AppUserRepo].
func (m *MySqlUserRepo) GetByIds(ids ...int) ([]*user.AppUser, error) {

	var users []*sqldto.UserSqlDto
	err := bun.NewSelectQuery(m.db).Where("id in (?)", bun.In(ids)).Model(&users).Scan(context.TODO())
	if err != nil {
		logger.Error().Err(err).Msg("UserRepo.GetByIds: bun query failed")
		return nil, err
	}

	appusers := funcs.Map(users, func(item *sqldto.UserSqlDto) *user.AppUser {
		return item.ToCoreUser()
	})

	return appusers, nil
}

// GetByPage implements [user.AppUserRepo].
func (m *MySqlUserRepo) GetByPage(filter *baserepo.PaginatedFilter[int]) (*baserepo.PaginatedResult[*user.AppUser, int], error) {

	var users []*user.AppUser
	var hasMore = true

	err := bun.NewSelectQuery(m.db).Where("id > ?", filter.Cursor).Limit(filter.Limit + 1).
		Model(&users).Scan(context.TODO())
	if err != nil {
		logger.Error().Err(err).Msg("UserRepo.GetByPage: failed with error")
		return nil, err
	}

	if len(users) != filter.Limit+1 {
		hasMore = false
		fmt.Println("Ok")
	} else {
		users = users[:filter.Limit]
	}

	return &baserepo.PaginatedResult[*user.AppUser, int]{
		Items:      users,
		HasMore:    hasMore,
		NextCursor: filter.Limit + 1,
	}, nil

}

// GetById implements user.AppUserRepo.
func (m *MySqlUserRepo) GetById(id int) (*user.AppUser, error) {
	var user sqldto.UserSqlDto
	err := bun.NewSelectQuery(m.db).Where("id = ?", id).Model(&user).Scan(context.TODO())
	return user.ToCoreUser(), err
}

func (m *MySqlUserRepo) GetByUserName(username string) (*user.AppUser, error) {

	var appuser user.AppUser

	f := filters.EQ("username", username)
	sqlVisitor := NewSqlVisitor()
	f.Accept(sqlVisitor)

	dsf := sq.Select("id", "email", "firstname", "lastname", "username").From("appuser").Where(sqlVisitor.String())
	fmt.Println(dsf.ToSql())
	row := dsf.RunWith(m.db).QueryRow()
	err := row.Scan(&appuser.Id, &appuser.Email, &appuser.Firstname, &appuser.Lastname, &appuser.UserName)
	if err != nil {
		return nil, err
	}
	return &appuser, nil
}

// GetPasswordHash implements user.AppUserRepo.
func (m *MySqlUserRepo) GetPasswordHash(id int) (string, error) {
	var str = "'"
	row := m.db.QueryRow("SELECT password FROM appuser WHERE id = ?", id)
	err := row.Scan(&str)

	return str, err
}

// Insert implements user.AppUserRepo.
func (m *MySqlUserRepo) Insert(val *user.AppUser) error {
	r, err := sq.Insert("appuser").Columns("UserName", "Password", "Firstname", "Lastname", "Email").
		Values(val.Email, val.Password, val.Firstname, val.Lastname, val.Email).
		RunWith(m.db).Exec()

	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return errors.New("could not fetch last insert id")
	}

	val.Id = int(id)

	return nil
}

// vscode interface impl
// var s user.AppUserRepo = &MySqlUserRepo{}
