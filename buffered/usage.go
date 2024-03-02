package buffered

import (
	"context"
)

type UsageState uint64

const (
	UsageUnknown UsageState = iota
	UsageL                  = iota
	UsageM                  = iota
	UsageH                  = iota
)

type UsageSource interface {
	Usage(context.Context) (UsageState, error)
}
