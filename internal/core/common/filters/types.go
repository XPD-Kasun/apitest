package filters

type And struct {
	Filters []Filter
}

func (a And) Accept(v Visitor) error {
	return v.VisitAnd(a)
}

type Or struct {
	Filters []Filter
}

func (or Or) Accept(v Visitor) error {
	return v.VisitOr(or)
}

type Not struct {
	Child Filter
}

func (not Not) Accept(v Visitor) error {
	return v.VisitNot(not)
}

type Eq struct {
	Field string
	Value any
}

func (op Eq) Accept(v Visitor) error {
	return v.VisitEq(op)
}

type Nq struct {
	Field string
	Value any
}

func (op Nq) Accept(v Visitor) error {
	return v.VisitNq(op)
}

type Gt struct {
	Field string
	Value any
}

func (op Gt) Accept(v Visitor) error {
	return v.VisitGt(op)
}

type Gte struct {
	Field string
	Value any
}

func (op Gte) Accept(v Visitor) error {
	return v.VisitGte(op)

}

type Lt struct {
	Field string
	Value any
}

func (op Lt) Accept(v Visitor) error {
	return v.VisitLt(op)
}

type Lte struct {
	Field string
	Value any
}

func (op Lte) Accept(v Visitor) error {
	return v.VisitLte(op)
}

type Like struct {
	Field string
	Value string
}

func (op Like) Accept(v Visitor) error {
	return v.VisitLike(op)
}

type In struct {
	Field  string
	Values []any
}

func (op In) Accept(v Visitor) error {
	return v.VisitIn(op)
}

type Between struct {
	Field string
	From  any
	To    any
}

func (op Between) Accept(v Visitor) error {
	return v.VisitBetween(op)
}
