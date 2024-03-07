package pair

func ErrMap[R, S any](p Pair[error, R], f func(R) S) Pair[error, S] {
	switch p.Left {
	case nil:
		return Right[error](f(p.Right))
	default:
		return Pair[error, S]{Left: p.Left}
	}
}

func ErrAndThen[R, S any](
	p Pair[error, R],
	f func(R) Pair[error, S],
) Pair[error, S] {
	switch p.Left {
	case nil:
		return f(p.Right)
	default:
		return Pair[error, S]{Left: p.Left}
	}
}
