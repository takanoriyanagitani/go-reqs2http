package buffered

import (
	"time"
)

const (
	WaitStaticDefault time.Duration = 1000 * time.Second
)

type WaitHint interface {
	Hint(UsageState, ChangeState) time.Duration
}

type WaitHintFn func(UsageState, ChangeState) time.Duration

func (f WaitHintFn) Hint(u UsageState, c ChangeState) time.Duration {
	return f(u, c)
}

func (f WaitHintFn) AsIf() WaitHint { return f }

func WaitHintFnStaticNew(s time.Duration) WaitHintFn {
	return WaitHintFn(func(_ UsageState, _ ChangeState) time.Duration {
		return s
	})
}

var WaitHintFnStaticDefault WaitHintFn = WaitHintFnStaticNew(WaitStaticDefault)
var WaitHintStaticDefault WaitHint = WaitHintFnStaticDefault.AsIf()
