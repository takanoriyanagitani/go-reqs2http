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

func Slice2Chan[T, U any](
	ctx context.Context,
	src []T,
	mkch func() chan U,
	conv func(T) U,
) <-chan U {
	var ch chan U = mkch()
	go func() {
		defer close(ch)

		for _, t := range src {
			var u U = conv(t)
			select {
			case <-ctx.Done():
				return
			default:
				ch <- u
			}
		}
	}()
	return ch
}

func ReqSrcChanFnFromSlice(s []*rhp.Request) ReqSrcChanFn {
	return func(ctx context.Context) <-chan RequestResult {
		return Slice2Chan(
			ctx,
			s,
			func() chan RequestResult { return make(chan RequestResult) },
			func(q *rhp.Request) RequestResult { return RequestOk(q) },
		)
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

//revive:disable:cognitive-complexity
func AutoStopChan[T any](
	ctx context.Context,
	mkch func() chan T,
	next func(context.Context) (T, error),
	stop func(error) bool,
) <-chan T {
	var ch chan T = mkch()
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			t, e := next(ctx)
			if stop(e) {
				return
			}
			ch <- t
		}
	}()
	return ch
}

//revive:enable:cognitive-complexity

func (f RequestSrcFn) ToChan(bufSz int) RequestSourceCh {
	return ReqSrcChanFn(func(ctx context.Context) <-chan RequestResult {
		var src RequestSource = f.AsIf()

		return AutoStopChan(
			ctx,
			func() chan RequestResult {
				return make(chan RequestResult, bufSz)
			},
			func(ctx context.Context) (RequestResult, error) {
				req, e := src.Next(ctx)
				return RequestResult{Left: e, Right: req}, e
			},
			func(e error) (stop bool) { return errors.Is(e, ErrNoMoreData) },
		)
	})
}

func RequestSrcFnFromSlice(s []*rhp.Request) RequestSrcFn {
	var ix int = 0
	var sz int = len(s)
	return RequestSrcFn(func(ctx context.Context) (*rhp.Request, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			break
		}

		if ix < sz {
			var req *rhp.Request = s[ix]
			ix++
			return req, nil
		}

		return nil, ErrNoMoreData
	})
}

var RequestSrcFnEmpty RequestSrcFn = RequestSrcFnFromSlice(nil)
var RequestSourceEmpty RequestSource = RequestSrcFnEmpty.AsIf()
