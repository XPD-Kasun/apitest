package sql

import (
	"apitest/internal/core/common/baserepo"
	"apitest/internal/core/common/filters"
	"apitest/internal/core/user"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

func NewMySqlUserRepo(db *sql.DB) (*MySqlUserRepo, error) {
	return &MySqlUserRepo{db: db}, nil
}

type MySqlUserRepo struct {
	db *sql.DB
}

func bindUser(row *sql.Row) (*user.AppUser, error) {

	u := user.AppUser{}
	err := row.Scan(&u.Id, &u.UserName, &u.Password, &u.Firstname, &u.Lastname, &u.Email)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func bindUsers(rows *sql.Rows, expectedCount int) (users []*user.AppUser, total int, err error) {

	users = make([]*user.AppUser, 0, 10)

	for rows.Next() {
		u := user.AppUser{}
		total++
		err := rows.Scan(&u.Id, &u.UserName, &u.Password, &u.Firstname, &u.Lastname, &u.Email)
		if err != nil && expectedCount != total {
			log.Error().Err(err).Msg("at bind users")
			continue
		}
		users = append(users, &u)
	}

	return
}

// GetByIds implements [user.AppUserRepo].
func (m *MySqlUserRepo) GetByIds(ids ...int) ([]*user.AppUser, error) {

	q := sq.Select("Id", "UserName", "Password", "Firstname", "Lastname", "Email").From("appuser").
		Where(sq.Eq{"Id": ids})

	sql, _, _ := q.ToSql()
	log.Debug().Str("sql", sql).Msg("GetByIds in MysqlUserRepo")

	rows, err := q.RunWith(m.db).Query()
	if err != nil {
		return nil, err
	}

	users, _, err := bindUsers(rows, len(ids))
	return users, err
}

// GetByPage implements [user.AppUserRepo].
func (m *MySqlUserRepo) GetByPage(filter *baserepo.PaginatedFilter[int]) (*baserepo.PaginatedResult[*user.AppUser, int], error) {

	q := sq.Select("Id", "UserName", "Password", "Firstname", "Lastname", "Email").From("appuser").
		Where(sq.GtOrEq{"Id": filter.Cursor}).OrderBy("Id asc").Limit(uint64(filter.Limit) + 1)

	sql, args := q.MustSql()

	log.Debug().Interface("args", args).Str("sql", sql).Int("cursor", filter.Cursor).
		Int("limit", filter.Limit).Msg("GetByPage invoked")

	rows, err := q.RunWith(m.db).Query()

	if err != nil {
		return nil, err
	}

	u, total, err := bindUsers(rows, filter.Limit+1)
	if err != nil {
		return nil, err
	}

	hasMore, itemLimit := false, len(u)
	if total == filter.Limit+1 {
		hasMore = true
		itemLimit--
	}

	return &baserepo.PaginatedResult[*user.AppUser, int]{
		Items:      u[:itemLimit],
		HasMore:    hasMore,
		NextCursor: filter.Limit + 1,
	}, nil

}

// GetById implements user.AppUserRepo.
func (m *MySqlUserRepo) GetById(id int) (*user.AppUser, error) {
	row := m.db.QueryRow("SELECT * FROM appuser WHERE id = $1", id)
	appUser := &user.AppUser{}
	err := row.Scan(&appUser.Id, &appUser.Email, &appUser.Firstname, &appUser.Lastname, &appUser.UserName, &appUser.Password)

	return appUser, err
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
