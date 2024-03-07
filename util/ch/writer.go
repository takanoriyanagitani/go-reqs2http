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

func WriteToChFromBuilder[T any](
	builder func(context.Context) (WriteToCh[T], error),
) WriteToCh[T] {
	return WriteToCh[T](func(
		ctx context.Context,
		dst chan<- pair.Pair[error, T],
	) error {
		w2c, e := builder(ctx)
		if nil == e {
			return w2c(ctx, dst)
		}
		return e
	})
}

type Write2ChConverter[S, T any] func(context.Context, S) (WriteToCh[T], error)

func (f Write2ChConverter[S, T]) ToConverted(src S) WriteToCh[T] {
	return WriteToCh[T](func(
		ctx context.Context,
		ch chan<- pair.Pair[error, T],
	) error {
		w2c, e := f(ctx, src)
		if nil == e {
			return w2c(ctx, ch)
		}
		return e
	})
}

func WriteMany[S, T any](
	ctx context.Context,
	sources <-chan pair.Pair[error, S],
	target chan<- pair.Pair[error, T],
	s2wt Write2ChConverter[S, T],
) error {
	return TryForEach(
		ctx,
		sources,
		func(src S) error {
			w2c, e := s2wt(ctx, src)
			if nil != e {
				return e
			}
			return w2c(ctx, target)
		},
	)
}

func (f Write2ChConverter[S, T]) WriteMany(
	ctx context.Context,
	sources <-chan pair.Pair[error, S],
	target chan<- pair.Pair[error, T],
) error {
	return WriteMany(ctx, sources, target, f)
}
