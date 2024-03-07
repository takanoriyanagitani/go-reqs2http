package chmap

import (
	"context"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"
)

func MapErr[T, U any](
	ctx context.Context,
	src <-chan pair.Pair[error, T],
	mapper func(context.Context, T) (U, error),
) <-chan pair.Pair[error, U] {
	ret := make(chan pair.Pair[error, U])
	go func() {
		defer close(ret)

		e := uch.TryForEach(
			ctx,
			src,
			func(t T) error {
				u, e := mapper(ctx, t)
				if nil != e {
					return e
				}
				ret <- pair.Right[error](u)
				return nil
			},
		)

		if nil != e {
			ret <- pair.Pair[error, U]{Left: e}
		}
	}()
	return ret
}

type ConvErr[T, U any] func(context.Context, T) (U, error)

func (f ConvErr[T, U]) MapErr(
	ctx context.Context,
	src <-chan pair.Pair[error, T],
) <-chan pair.Pair[error, U] {
	return MapErr(
		ctx,
		src,
		f,
	)
}
