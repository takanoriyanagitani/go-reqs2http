package arr

func All[T any](s []T, f func(T) bool) bool {
	return Fold(
		s,
		true,
		func(state bool, next T) bool { return state && f(next) },
	)
}
