package util

func Select[T any](ng T, ok T, cond bool) T {
	switch cond {
	case true:
		return ok
	default:
		return ng
	}
}
