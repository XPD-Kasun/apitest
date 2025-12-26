package filters

func EQ(field string, value any) Filter {
	f := Eq{Field: field, Value: value}
	return f
}

func NQ(field string, value any) Filter {
	f := Nq{Field: field, Value: value}
	return f
}

func GT(field string, value any) Filter {
	f := Gt{Field: field, Value: value}
	return f
}

func GTE(field string, value any) Filter {
	f := Gte{Field: field, Value: value}
	return f
}

func LT(field string, value any) Filter {
	f := Lt{Field: field, Value: value}
	return f
}

func LTE(field string, value any) Filter {
	f := Lte{Field: field, Value: value}
	return f
}

func AND(filters ...Filter) Filter {
	f := And{Filters: filters}
	return f
}

func OR(filters ...Filter) Filter {
	f := Or{Filters: filters}
	return f
}

func NOT(expr Filter) Filter {
	f := Not{Child: expr}
	return f
}

func LIKE(field string, value string) Filter {
	f := Like{Field: field, Value: value}
	return f
}

func IN(field string, args ...any) Filter {
	f := In{
		Field:  field,
		Values: args,
	}
	return f
}

func BETWEEN(field string, from, to any) Filter {
	f := Between{Field: field, From: from, To: to}
	return f
}
