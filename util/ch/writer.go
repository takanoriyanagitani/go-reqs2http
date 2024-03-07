package ch

import (
	"context"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
)

type WriterToChan[T any] interface {
	WriteToChan(ctx context.Context, ch chan<- T) error
}

type WriteToCh[T any] func(context.Context, chan<- pair.Pair[error, T]) error

func (f WriteToCh[T]) WriteToChan(
	ctx context.Context,
	ch chan<- pair.Pair[error, T],
) error {
	return f(ctx, ch)
}

func (f WriteToCh[T]) AsIf() WriterToChan[pair.Pair[error, T]] { return f }
func (f WriteToCh[T]) ToChan(ctx context.Context) <-chan pair.Pair[error, T] {
	c := make(chan pair.Pair[error, T])
	go func() {
		defer close(c)
		e := f.WriteToChan(ctx, c)
		if nil == e {
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			c <- pair.Pair[error, T]{Left: e}
		}
	}()
	return c
}
func (f WriteToCh[T]) Fold(
	ctx context.Context,
	init pair.Pair[error, T],
	reducer func(
		state pair.Pair[error, T],
		next pair.Pair[error, T],
	) pair.Pair[error, T],
) pair.Pair[error, T] {
	var ch <-chan pair.Pair[error, T] = f.ToChan(ctx)
	return Fold(
		ctx,
		init,
		ch,
		reducer,
	)
}

func (f WriteToCh[T]) TryFold(
	ctx context.Context,
	init T,
	reducer func(state T, next T) pair.Pair[error, T],
) pair.Pair[error, T] {
	var ch <-chan pair.Pair[error, T] = f.ToChan(ctx)
	return TryFold(
		ctx,
		init,
		ch,
		reducer,
	)
}

func WriteToChFromSlice[T any](s []T) WriteToCh[T] {
	return WriteToCh[T](func(
		ctx context.Context,
		dst chan<- pair.Pair[error, T],
	) error {
		for _, item := range s {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				dst <- pair.Right[error](item)
			}
		}
		return nil
	})
}
