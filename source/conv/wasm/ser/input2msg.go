package wser

import (
	"context"

	cser "github.com/takanoriyanagitani/go-reqs2http/source/conv/ser"

	wcnv "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm"
)

type Input2msg[M any] struct {
	Converter wcnv.Converter
	Bytes2Msg cser.ConvertFn[M]
}

func (c Input2msg[M]) ToMessage(
	ctx context.Context,
	input []byte,
	buf M,
) error {
	converted, e := c.Converter.Convert(ctx, input)
	if nil != e {
		return e
	}
	return c.Bytes2Msg(converted, buf)
}
