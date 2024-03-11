package arr

func TryForEach[T any](s []T, f func(T) error) error {
	for _, item := range s {
		e := f(item)
		if nil != e {
			return e
		}
	}
	return nil
}
