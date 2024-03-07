package util

import (
	"context"
)

func ComposeCtx[T, U, V any](
	f func(context.Context, T) (U, error),
	g func(context.Context, U) (V, error),
) func(context.Context, T) (V, error) {
	return func(ctx context.Context, t T) (v V, e error) {
		u, e := f(ctx, t)
		return Select(
			func() (V, error) { return v, e },
			func() (V, error) { return g(ctx, u) },
			nil == e,
		)()
	}
}

func CtxIgnore[T, U any](
	f func(T) (U, error),
) func(context.Context, T) (U, error) {
	return func(_ context.Context, t T) (U, error) { return f(t) }
}
