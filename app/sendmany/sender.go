package sendmany

import (
	"context"
	"errors"
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

func (m ManySender) getWait(ctx context.Context) (time.Duration, error) {
	usg, eu := m.sender.Usage(ctx)
	chg, ec := m.sender.Change(ctx)
	e := errors.Join(eu, ec)
	if nil != e {
		return 1000 * time.Second, e
	}
	return m.waiter.Hint(usg, chg), nil
}

//revive:disable:cognitive-complexity
func ProcessChan[T any](
	ctx context.Context,
	ch <-chan T,
	onData func(T) error,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case dat, ok := <-ch:
			if !ok {
				return nil
			}

			e := onData(dat)
			if nil != e {
				return e
			}
		}
	}
}

//revive:enable:cognitive-complexity

func (m ManySender) SendAll(ctx context.Context, bufSz int) error {
	var sf src.RequestSrcFn = src.RequestSrcFn(m.source.Next)
	var sc src.RequestSourceCh = sf.ToChan(bufSz)
	var reqs <-chan src.RequestResult = sc.GetRequests(ctx)
	return ProcessChan(
		ctx,
		reqs,
		func(rslt src.RequestResult) error {
			var err error = rslt.Left
			if nil != err {
				return err
			}

			var req *rhp.Request = rslt.Right
			e := m.sender.Push(ctx, req)
			if nil != e {
				return e
			}

			wait, e := m.getWait(ctx)
			if nil != e {
				return e
			}

			time.Sleep(wait)

			return nil
		},
	)
}

func (m ManySender) WithSource(s src.RequestSource) ManySender {
	m.source = s
	return m
}

func (m ManySender) WithSender(s buf.Sender) ManySender {
	m.sender = s
	return m
}

func (m ManySender) WithWaiter(w buf.WaitHint) ManySender {
	m.waiter = w
	return m
}

var ManySenderNopDefault ManySender = ManySender{
	source: src.RequestSourceEmpty,
	sender: buf.SenderNop,
	waiter: buf.WaitHintStaticDefault,
}
