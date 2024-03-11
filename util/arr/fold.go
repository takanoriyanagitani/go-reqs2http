package arr

func Fold[T, S any](s []T, init S, reducer func(S, T) S) S {
	var state S = init
	for _, next := range s {
		state = reducer(state, next)
	}
	return state
}
