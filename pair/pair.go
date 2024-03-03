package pair

type Pair[L, R any] struct {
	Left  L
	Right R
}

func Right[L, R any](right R) Pair[L, R] {
	return Pair[L, R]{Right: right}
}
