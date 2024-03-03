package buffered

import (
	"time"
)

type WaitHint interface {
	Hint(UsageState, ChangeState) time.Duration
}

type WaitHintFn func(UsageState, ChangeState) time.Duration

func (f WaitHintFn) Hint(u UsageState, c ChangeState) time.Duration {
	return f(u, c)
}

func (f WaitHintFn) AsIf() WaitHint { return f }
