package buffered

import (
	"time"
)

type WaitHint interface {
	Hint(UsageState, ChangeState) time.Duration
}
