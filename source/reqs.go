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

type RequestResult pair.Pair[error, *rhp.Request]

func RequestOk(r *rhp.Request) RequestResult {
	return RequestResult(pair.Right[error](r))
}

type RequestSourceCh interface {
	GetRequests(context.Context) <-chan RequestResult
}

type ReqSrcChanFn func(context.Context) <-chan RequestResult

func (f ReqSrcChanFn) GetRequests(ctx context.Context) <-chan RequestResult {
	return f(ctx)
}
func (f ReqSrcChanFn) AsIf() RequestSourceCh { return f }
func (f ReqSrcChanFn) GetAll(ctx context.Context) (pairs []RequestResult) {
	var pch <-chan RequestResult = f(ctx)

	for pair := range pch {
		pairs = append(pairs, pair)
	}
	return
}

func ReqSrcChanFnFromSlice(s []*rhp.Request) ReqSrcChanFn {
	return func(ctx context.Context) <-chan RequestResult {
		ch := make(chan RequestResult)
		go func() {
			defer close(ch)

			for _, req := range s {
				select {
				case <-ctx.Done():
					return
				default:
					ch <- RequestOk(req)
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

func (f RequestSrcFn) Next(ctx context.Context) (*rhp.Request, error) {
	return f(ctx)
}

func (f RequestSrcFn) AsIf() RequestSource { return f }

func (f RequestSrcFn) ToChan(bufSz int) RequestSourceCh {
	return ReqSrcChanFn(func(ctx context.Context) <-chan RequestResult {
		ret := make(chan RequestResult, bufSz)
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
				ret <- RequestResult{Left: e, Right: req}
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
