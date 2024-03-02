package waithint

import (
	"time"

	buf "github.com/takanoriyanagitani/go-reqs2http/buffered"
)

type SimpleHint struct {
	Short  time.Duration
	Normal time.Duration
	Long   time.Duration
}

func (s SimpleHint) Hint(u buf.UsageState, _ buf.ChangeState) time.Duration {
	switch u {
	case buf.UsageL:
		return s.Short
	case buf.UsageM:
		return s.Normal
	case buf.UsageH:
		return s.Long
	default:
		return s.Long
	}
}

func (s SimpleHint) AsIf() buf.WaitHint { return s }
