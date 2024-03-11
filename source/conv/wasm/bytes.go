package wasm

import (
	"context"
)

type Converter interface {
	Convert(ctx context.Context, input []byte) (output []byte, e error)
}

type ConvertFn func(context.Context, []byte) ([]byte, error)

func (f ConvertFn) Convert(ctx context.Context, i []byte) ([]byte, error) {
	return f(ctx, i)
}

func (f ConvertFn) AsIf() Converter { return f }
