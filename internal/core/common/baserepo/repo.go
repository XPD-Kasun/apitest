package baserepo

import (
	"apitest/internal/core/common/filters"
)

type SingleIdGetter[TVal any, TKey any] interface {
	GetById(id TKey) (TVal, error)
}

type MultiIdGetter[TVal any, TKey any] interface {
	GetByIds(ids ...TKey) ([]TVal, error)
}

type FilterGetter[TVal any] interface {
	Get(filters ...filters.Filter) ([]TVal, error)
}

type Inserter[TVal any] interface {
	Insert(val TVal) error
}

type Deleter[TVal any] interface {
	Delete(val TVal) error
}

type Updater[TVal any] interface {
	Update(val TVal) error
}

//pagination

type PaginatedFilter[TKey any] struct {
	Cursor TKey
	Limit  int
}

type PaginatedResult[TVal any, TKey any] struct {
	Items      []TVal
	HasMore    bool
	NextCursor TKey
}

type PaginatedGetter[TVal any, TKey any] interface {
	GetByPage(filter *PaginatedFilter[TKey]) (*PaginatedResult[TVal, TKey], error)
}
