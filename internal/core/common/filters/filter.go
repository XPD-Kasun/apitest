package filters

type Filter interface {
	Accept(Visitor) error
}

type Visitor interface {
	VisitAnd(And) error
	VisitOr(Or) error
	VisitNot(Not) error
	VisitEq(Eq) error
	VisitNq(Nq) error
	VisitGt(Gt) error
	VisitGte(Gte) error
	VisitLt(Lt) error
	VisitLte(Lte) error
	VisitLike(Like) error
	VisitIn(In) error
	VisitBetween(Between) error
}
