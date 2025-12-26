//go:build !dev

package op

type Operator int

const (
	Eq Operator = iota
	Nq
	Gt
	Gte
	Lt
	Lte
	And
	Or
	Not
	Like
	In
	Between
)
