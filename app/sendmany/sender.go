package sendmany

import (
	"context"
	"errors"
	"time"

	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"

	buf "github.com/takanoriyanagitani/go-reqs2http/buffered"
	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
	src "github.com/takanoriyanagitani/go-reqs2http/source"
)

type ManySender struct {
	source src.RequestSource
	sender buf.Sender
	waiter buf.WaitHint
	chsrc  src.RequestSourceCh
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

func ProcessChan[T any](
	ctx context.Context,
	ch <-chan T,
	onData func(T) error,
) error {
	return uch.ProcessChan(ctx, ch, onData)
}

func (m ManySender) SendAll(ctx context.Context, _ int) error {
	var sc src.RequestSourceCh = m.chsrc
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
	m.chsrc = src.RequestSrcFn(s.Next).ToChan(0)
	return m
}

func (m ManySender) WithSrcCh(s src.RequestSourceCh) ManySender {
	m.chsrc = s
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
	chsrc:  src.RequestSrcFn(src.RequestSourceEmpty.Next).ToChan(0),
}
