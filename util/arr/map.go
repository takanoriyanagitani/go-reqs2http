package arr

func Map[T, U any](s []T, mapper func(T) U) []U {
	var ret []U = make([]U, 0, len(s))
	for _, t := range s {
		var u U = mapper(t)
		ret = append(ret, u)
	}
	return ret
}

func MapErr[T, U any](s []T, mapper func(T) (U, error)) ([]U, error) {
	var ret []U = make([]U, 0, len(s))
	for _, t := range s {
		u, e := mapper(t)
		if nil != e {
			return nil, e
		}
		ret = append(ret, u)
	}
	return ret, nil
}
