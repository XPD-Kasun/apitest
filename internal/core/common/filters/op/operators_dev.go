//go:build dev

package op

type Operator string

const (
	Eq      Operator = "="
	Nq      Operator = "!="
	Gt      Operator = ">"
	Gte     Operator = ">="
	Lt      Operator = "<"
	Lte     Operator = "<="
	And     Operator = "&"
	Or      Operator = "|"
	Not     Operator = "!"
	Like    Operator = "LIKE "
	In      Operator = "IN "
	Between Operator = "BETWEEN "
)
