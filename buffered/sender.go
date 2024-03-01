package buffered

import (
	"context"
	"time"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type UsageState uint64

const (
	UsageUnknown UsageState = iota
	UsageL                  = iota
	UsageM                  = iota
	UsageH                  = iota
)

type ChangeState uint64

const (
	ChangeUnknown ChangeState = iota
	ChangeRise                = iota
	ChangeFlat                = iota
	ChangeFall                = iota
)

type UsageSource interface {
	Usage(context.Context) (UsageState, error)
}

type ChangeSource interface {
	Change(context.Context) (ChangeState, error)
}

type Buffer interface {
	Push(context.Context, *rhp.Request) error
}

type BufferedSender interface {
	UsageSource
	ChangeSource
	Buffer
}

type WaitHint interface {
	Hint(UsageState, ChangeState) time.Duration
}
