package wser

import (
	"context"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"

	cser "github.com/takanoriyanagitani/go-reqs2http/source/conv/ser"

	wcnv "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm"
)

type Input2requests[M any] struct {
	In2out    wcnv.Converter
	Bytes2msg cser.ConvertFn[M]
	Msg2reqs  cser.Message2requests[M]
}

func (c Input2requests[M]) Input2chan(
	ctx context.Context,
	input []byte,
	buf M,
	dst chan<- pair.Pair[error, *rhp.Request],
) error {
	serialized, e := c.In2out.Convert(ctx, input)
	if nil != e {
		return e
	}
	return c.Msg2reqs.Bytes2Chan(
		ctx,
		serialized,
		c.Bytes2msg,
		buf,
		dst,
	)
}
