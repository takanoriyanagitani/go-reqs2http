package ch

import (
	"context"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
)

//revive:disable:cognitive-complexity
func ProcessChan[T any](
	ctx context.Context,
	ch <-chan T,
	f func(T) error,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case next, ok := <-ch:
			if !ok {
				return nil
			}
			e := f(next)
			if nil != e {
				return e
			}
		}
	}
}

//revive:enable:cognitive-complexity

func ForEach[T any](
	ctx context.Context,
	ch <-chan T,
	f func(T),
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case next, ok := <-ch:
			if !ok {
				return nil
			}
			f(next)
		}
	}
}

func TryForEach[T any](
	ctx context.Context,
	ch <-chan pair.Pair[error, T],
	f func(T) error,
) error {
	return ProcessChan(
		ctx,
		ch,
		func(p pair.Pair[error, T]) error {
			if nil != p.Left {
				return p.Left
			}
			return f(p.Right)
		},
	)
}
