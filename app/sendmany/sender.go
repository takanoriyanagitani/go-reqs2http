package sendmany

import (
	"context"
	"time"

	buf "github.com/takanoriyanagitani/go-reqs2http/buffered"
	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
	src "github.com/takanoriyanagitani/go-reqs2http/source"
)

type ManySender struct {
	source src.RequestSource
	sender buf.Sender
	waiter buf.WaitHint
}

func (m ManySender) SendAll(ctx context.Context, bufSz int) error {
	var sf src.RequestSrcFn = src.RequestSrcFn(m.source.Next)
	var sc src.RequestSourceCh = sf.ToChan(bufSz)
	var reqs <-chan src.RequestResult = sc.GetRequests(ctx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p, ok := <-reqs:
			var nodata bool = !ok
			if nodata {
				return nil
			}

			var err error = p.Left
			if nil != err {
				return err
			}

			var req *rhp.Request = p.Right
			e := m.sender.Push(ctx, req)
			if nil != e {
				return e
			}

			usage, e := m.sender.Usage(ctx)
			if nil != e {
				return e
			}

			change, e := m.sender.Change(ctx)
			if nil != e {
				return e
			}

			var wait time.Duration = m.waiter.Hint(usage, change)

			time.Sleep(wait)
		}
	}
}
