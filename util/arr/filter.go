package arr

func FilterSlow[T any](s []T, f func(T) (keep bool)) []T {
	var ret []T
	for _, next := range s {
		var keep bool = f(next)
		if keep {
			ret = append(ret, next)
		}
	}
	return ret
}
