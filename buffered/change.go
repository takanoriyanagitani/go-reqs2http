package buffered

import (
	"context"
)

type ChangeState uint64

const (
	ChangeUnknown ChangeState = iota
	ChangeRise                = iota
	ChangeFlat                = iota
	ChangeFall                = iota
)

type ChangeSource interface {
	Change(context.Context) (ChangeState, error)
}

type ChangeSrcFn func(context.Context) (ChangeState, error)

func (f ChangeSrcFn) Change(c context.Context) (ChangeState, error) {
	return f(c)
}

func (f ChangeSrcFn) AsIf() ChangeSource { return f }

func ChangeSrcFnStaticNew(s ChangeState) ChangeSrcFn {
	return ChangeSrcFn(func(_ context.Context) (ChangeState, error) {
		return s, nil
	})
}

var ChangeSrcFnUnknown ChangeSrcFn = ChangeSrcFnStaticNew(ChangeUnknown)
var ChangeSourceUnknown ChangeSource = ChangeSrcFnUnknown.AsIf()
