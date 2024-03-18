package reqs2http

import (
	"context"
	"errors"
	"net/http"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

var (
	ErrInvalidRequest error = errors.New("invalid request")
)

// Sender sends an [http.Request].
type Sender interface {
	// Send tries to get an [http.Response] by sending the request.
	Send(context.Context, *http.Request) (*http.Response, error)
}

type SendFn func(context.Context, *http.Request) (*http.Response, error)

func (f SendFn) Send(
	c context.Context, q *http.Request,
) (*http.Response, error) {
	return f(c, q)
}

func (f SendFn) AsIf() Sender { return f }

func SenderNew(client *http.Client) Sender {
	return SendFn(func(
		ctx context.Context, req *http.Request,
	) (*http.Response, error) {
		var neo *http.Request = req.WithContext(ctx)
		return client.Do(neo)
	})
}

// SenderDefault uses [http.DefaultClient] to send an [http.Request].
var SenderDefault Sender = SenderNew(http.DefaultClient)

// A sender to send unconverted requests.
type RawSender interface {
	Send(context.Context, *rhp.Request) (*http.Response, error)
}

type RawSendFn func(context.Context, *rhp.Request) (*http.Response, error)

func (f RawSendFn) Send(
	c context.Context, q *rhp.Request,
) (*http.Response, error) {
	return f(c, q)
}

func (f RawSendFn) AsIf() RawSender { return f }

func RawSenderNew(sender Sender) func(RequestConverter) RawSender {
	return func(conv RequestConverter) RawSender {
		return RawSendFn(func(
			ctx context.Context, req *rhp.Request,
		) (*http.Response, error) {
			q, e := conv.Convert(req)
			if nil != e {
				return nil, errors.Join(ErrInvalidRequest, e)
			}
			return sender.Send(ctx, q)
		})
	}
}

var RawSenderDefault RawSender = RawSenderNew(SenderDefault)(RequestConvDefault)
