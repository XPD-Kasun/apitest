package sql

import (
	"apitest/internal/core/common/filters"
	"strings"
)

type SqlVisitor struct {
	buf  *strings.Builder
	args []string
}

// VisitBetween implements filters.Visitor.
func (v *SqlVisitor) VisitBetween(a filters.Between) error {
	w := v.buf
	w.WriteString(a.Field)
	w.WriteString(" BETWEEN $1 and $2")

	if arg1, err := Parse(a.From); err == nil {
		v.args = append(v.args, arg1)
	} else {
		return err
	}

	if arg2, err := Parse(a.To); err == nil {
		v.args = append(v.args, arg2)
	} else {
		return err
	}

	return nil
}

// VisitEq implements filters.Visitor.
func (v *SqlVisitor) VisitEq(op filters.Eq) error {

	return comparison(v.buf, op.Field, " = ", op.Value)

}

func comparison(w *strings.Builder, field string, op string, val any) error {

	w.WriteString(field)
	w.WriteString(op)

	if valx, err := Parse(val); err == nil {
		w.WriteString(valx)
	} else {
		return err
	}
	return nil
}

// VisitGt implements filters.Visitor.
func (v *SqlVisitor) VisitGt(op filters.Gt) error {
	return comparison(v.buf, op.Field, " > ", op.Value)
}

// VisitGte implements filters.Visitor.
func (v *SqlVisitor) VisitGte(op filters.Gte) error {
	return comparison(v.buf, op.Field, " >= ", op.Value)
}

// VisitIn implements filters.Visitor.
func (v *SqlVisitor) VisitIn(op filters.In) error {
	w := v.buf

	w.WriteString(op.Field)
	w.WriteString(" IN (")

	var s []string = make([]string, len(op.Values))

	for i, val := range op.Values {
		if v, err := Parse(val); err == nil {
			s[i] = v
		}
	}

	w.WriteString(strings.Join(s, ","))
	w.WriteString(")")

	return nil
}

// VisitLike implements filters.Visitor.
func (v *SqlVisitor) VisitLike(op filters.Like) error {
	panic("unimplemented")
}

// VisitLt implements filters.Visitor.
func (v *SqlVisitor) VisitLt(op filters.Lt) error {
	return comparison(v.buf, op.Field, " > ", op.Value)
}

// VisitLte implements filters.Visitor.
func (v *SqlVisitor) VisitLte(op filters.Lte) error {
	return comparison(v.buf, op.Field, " > ", op.Value)
}

// VisitNq implements filters.Visitor.
func (v *SqlVisitor) VisitNq(op filters.Nq) error {
	return comparison(v.buf, op.Field, " != ", op.Value)
}

func NewSqlVisitor() *SqlVisitor {
	return &SqlVisitor{
		buf: &strings.Builder{},
	}
}

func (v *SqlVisitor) VisitAnd(a filters.And) error {
	w := v.buf
	for i := 0; i < len(a.Filters); i++ {
		w.WriteString("(")
		f := a.Filters[i]
		err := f.Accept(v)
		if err != nil {
			return err
		}
		if i != (len(a.Filters) - 1) {
			w.WriteString(") AND ")
		} else {
			w.WriteString(")")
		}

	}
	return nil
}

func (v *SqlVisitor) VisitOr(or filters.Or) error {
	w := v.buf
	for i := 0; i < len(or.Filters); i++ {
		w.WriteString("(")
		f := or.Filters[i]
		err := f.Accept(v)
		if err != nil {
			return err
		}
		if i != (len(or.Filters) - 1) {
			w.WriteString(") OR ")
		} else {
			w.WriteString(")")
		}

	}
	return nil
}

func (v *SqlVisitor) VisitNot(not filters.Not) error {
	w := v.buf
	w.WriteString("NOT (")
	err := not.Child.Accept(v)
	if err != nil {
		return err
	}
	w.WriteString(")")
	return nil
}

func (v *SqlVisitor) String() string {
	return v.buf.String()
}
