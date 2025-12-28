package funcs

func Map[TSrc any, TDest any](ar []TSrc, fn func(TSrc) TDest) []TDest {
	var result []TDest = make([]TDest, len(ar))

	for i, item := range ar {
		result[i] = fn(item)
	}

	return result
}
