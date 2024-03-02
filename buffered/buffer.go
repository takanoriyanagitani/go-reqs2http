package buffered

import (
	"context"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type Buffer interface {
	Push(context.Context, *rhp.Request) error
}
