package cser

import (
	"context"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type Message2requests[M any] func(
	ctx context.Context,
	msg M,
	dst chan<- pair.Pair[error, *rhp.Request],
) error

func (f Message2requests[M]) Bytes2Chan(
	ctx context.Context,
	serialized []byte,
	converter ConvertFn[M],
	buf M,
	dst chan<- pair.Pair[error, *rhp.Request],
) error {
	e := converter(serialized, buf)
	if nil != e {
		return e
	}
	return f(ctx, buf, dst)
}
