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
