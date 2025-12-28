package sqldto

import "github.com/uptrace/bun"

type UserSqlDto struct {
	bun.BaseModel `bun:"table:appuser"`

	Id        int
	UserName  string
	Firstname string
	Lastname  string
	Email     string
}
