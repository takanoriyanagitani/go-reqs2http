package chmap

import (
	"context"

	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"
)

func Map[T, U any](
	ctx context.Context,
	src <-chan T,
	mapper func(context.Context, T) U,
) <-chan U {
	ret := make(chan U)
	go func() {
		defer close(ret)

		_ = uch.ForEach(
			ctx,
			src,
			func(t T) {
				var mapd U = mapper(ctx, t)
				ret <- mapd
			},
		)
	}()
	return ret
}
