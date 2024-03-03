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
func (f ReqSrcChanFn) GetAll(ctx context.Context) (pairs []pair.Pair[error, *rhp.Request]) {
	var pch <-chan pair.Pair[error, *rhp.Request] = f(ctx)

	for pair := range pch {
		pairs = append(pairs, pair)
	}
	return
}

func ReqSrcChanFnFromSlice(s []*rhp.Request) ReqSrcChanFn {
	return func(ctx context.Context) <-chan pair.Pair[error, *rhp.Request] {
		ch := make(chan pair.Pair[error, *rhp.Request])
		go func() {
			defer close(ch)

			for _, req := range s {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- pair.Right[error](req)
				}
			}
		}()
		return ch
	}
}

type RequestSource interface {
	// Next returns a next request object.
	// If there's no more objects, returns nil, ErrNoMoreData.
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
				if errors.Is(e, ErrNoMoreData) {
					return
				}
				ret <- pair.Pair[error, *rhp.Request]{Left: e, Right: req}
			}
		}()
		return ret
	})
}

func RequestSrcFnFromSlice(s []*rhp.Request) RequestSrcFn {
	var ix int = 0
	var sz int = len(s)
	return RequestSrcFn(func(ctx context.Context) (*rhp.Request, error) {
		if ix < sz {
			var req *rhp.Request = s[ix]
			ix += 1
			return req, nil
		}

		return nil, ErrNoMoreData
	})
}

var RequestSrcFnEmpty RequestSrcFn = RequestSrcFnFromSlice(nil)
