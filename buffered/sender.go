package buffered

import (
	"context"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type Sender interface {
	UsageSource
	ChangeSource
	Buffer
}

type SenderFn struct {
	UsageFn  UsageSrcFn
	ChangeFn ChangeSrcFn
	BufferFn BufFn
}

func (f SenderFn) Usage(ctx context.Context) (UsageState, error) {
	return f.UsageFn(ctx)
}

func (f SenderFn) Change(ctx context.Context) (ChangeState, error) {
	return f.ChangeFn(ctx)
}

func (f SenderFn) Push(ctx context.Context, q *rhp.Request) error {
	return f.BufferFn(ctx, q)
}

func (f SenderFn) AsIf() Sender { return f }

var SenderNop Sender = SenderFn{
	UsageFn:  UsageSrcFnUnknown,
	ChangeFn: ChangeSrcFnUnknown,
	BufferFn: BufFnNop,
}
