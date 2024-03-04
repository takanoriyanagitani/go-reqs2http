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

type UsageSrcFn func(context.Context) (UsageState, error)

func (f UsageSrcFn) Usage(ctx context.Context) (UsageState, error) {
	return f(ctx)
}

func (f UsageSrcFn) AsIf() UsageSource { return f }

func UsageSrcFnStaticNew(s UsageState) UsageSrcFn {
	return UsageSrcFn(func(_ context.Context) (UsageState, error) {
		return s, nil
	})
}

var UsageSrcFnUnknown UsageSrcFn = UsageSrcFnStaticNew(UsageUnknown)
var UsageSourceUnknown UsageSource = UsageSrcFnUnknown.AsIf()
