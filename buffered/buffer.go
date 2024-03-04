package buffered

import (
	"context"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type Buffer interface {
	Push(context.Context, *rhp.Request) error
}

type BufFn func(context.Context, *rhp.Request) error

func (f BufFn) Push(c context.Context, q *rhp.Request) error { return f(c, q) }
func (f BufFn) AsIf() Buffer                                 { return f }

var BufFnNop BufFn = BufFn(func(_ context.Context, _ *rhp.Request) error {
	return nil
})

var BufferNop Buffer = BufFnNop.AsIf()
