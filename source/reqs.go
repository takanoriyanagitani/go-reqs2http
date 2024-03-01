package source

import (
	"context"
	"errors"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

var (
	ErrNoMoreData error = errors.New("no more data")
)

type RequestSourceCh interface {
	GetRequests(context.Context) <-chan pair.Pair[error, *rhp.Request]
}

type ReqSrcChanFn func(context.Context) <-chan pair.Pair[error, *rhp.Request]

func (f ReqSrcChanFn) GetRequests(ctx context.Context) <-chan pair.Pair[error, *rhp.Request] {
	return f(ctx)
}
func (f ReqSrcChanFn) AsIf() RequestSourceCh { return f }

type RequestSource interface {
	Next(context.Context) (*rhp.Request, error)
}

type RequestSrcFn func(context.Context) (*rhp.Request, error)

func (f RequestSrcFn) Next(ctx context.Context) (*rhp.Request, error) { return f(ctx) }
func (f RequestSrcFn) AsIf() RequestSource                            { return f }

func (f RequestSrcFn) ToChan(bufSz int) RequestSourceCh {
	return ReqSrcChanFn(func(ctx context.Context) <-chan pair.Pair[error, *rhp.Request] {
		ret := make(chan pair.Pair[error, *rhp.Request], bufSz)
		var src RequestSource = f.AsIf()
		go func() {
			defer close(ret)
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				req, e := src.Next(ctx)
				if ErrNoMoreData == e {
					return
				}
				ret <- pair.Pair[error, *rhp.Request]{Left: e, Right: req}
			}
		}()
		return ret
	})
}
