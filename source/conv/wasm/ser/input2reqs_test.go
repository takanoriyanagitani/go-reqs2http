package wser_test

import (
	"context"
	"testing"

	"encoding/json"
	"strings"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"

	ua "github.com/takanoriyanagitani/go-reqs2http/util/arr"
	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"

	cser "github.com/takanoriyanagitani/go-reqs2http/source/conv/ser"
	wcnv "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm"
	wser "github.com/takanoriyanagitani/go-reqs2http/source/conv/wasm/ser"
)

func TestSer(t *testing.T) {
	t.Parallel()

	t.Run("Input2requests", func(t *testing.T) {
		t.Parallel()

		t.Run("Input2chan", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var i2o wcnv.Converter = wcnv.ConvertFn(
					func(_ context.Context, _ []byte) (output []byte, e error) {
						return
					},
				).AsIf()

				var ser2buf cser.ConvertFn[[]*rhp.Request] = func(
					_ []byte,
					buf []*rhp.Request,
				) error {
					return nil
				}

				var msg2reqs cser.Message2requests[[]*rhp.Request] = func(
					ctx context.Context,
					msg []*rhp.Request,
					dst chan<- pair.Pair[error, *rhp.Request],
				) error {
					for _, m := range msg {
						select {
						case <-ctx.Done():
							return ctx.Err()
						default:
						}
						dst <- pair.Right[error](m)
					}
					return nil
				}

				i2r := wser.Input2requests[[]*rhp.Request]{
					In2out:    i2o,
					Bytes2msg: ser2buf,
					Msg2reqs:  msg2reqs,
				}

				reqs := make(chan pair.Pair[error, *rhp.Request])

				go func() {
					defer close(reqs)

					var buf []*rhp.Request
					e := i2r.Input2chan(
						context.Background(),
						nil,
						buf,
						reqs,
					)
					if nil != e {
						panic(e)
					}
				}()

				var cnt pair.Pair[error, int] = uch.TryFold(
					context.Background(),
					0,
					reqs,
					func(state int, _ *rhp.Request) pair.Pair[error, int] {
						return pair.Right[error](state + 1)
					},
				)
				t.Run("no chan err", assertNil(cnt.Left))
				t.Run("no items", assertEqual(cnt.Right, 0))
			})

			t.Run("dummy", func(t *testing.T) {
				t.Parallel()

				var i2o wcnv.Converter = wcnv.ConvertFn(
					func(_ context.Context, _ []byte) (output []byte, e error) {
						return
					},
				).AsIf()

				var ser2buf cser.ConvertFn[[]*rhp.Request] = func(
					_ []byte,
					buf []*rhp.Request,
				) error {
					buf[0] = &rhp.Request{}
					return nil
				}

				var msg2reqs cser.Message2requests[[]*rhp.Request] = func(
					ctx context.Context,
					msg []*rhp.Request,
					dst chan<- pair.Pair[error, *rhp.Request],
				) error {
					for _, m := range msg {
						select {
						case <-ctx.Done():
							return ctx.Err()
						default:
						}
						dst <- pair.Right[error](m)
					}
					return nil
				}

				i2r := wser.Input2requests[[]*rhp.Request]{
					In2out:    i2o,
					Bytes2msg: ser2buf,
					Msg2reqs:  msg2reqs,
				}

				reqs := make(chan pair.Pair[error, *rhp.Request])

				go func() {
					defer close(reqs)
					var buf []*rhp.Request = []*rhp.Request{
						{},
					}

					e := i2r.Input2chan(
						context.Background(),
						nil,
						buf,
						reqs,
					)
					t.Run("no err", assertNil(e))
				}()

				var cnt pair.Pair[error, int] = uch.TryFold(
					context.Background(),
					0,
					reqs,
					func(state int, _ *rhp.Request) pair.Pair[error, int] {
						return pair.Right[error](state + 1)
					},
				)
				t.Run("no chan err", assertNil(cnt.Left))
				t.Run("single item", assertEqual(cnt.Right, 1))

			})

			t.Run("jarr2strings", func(t *testing.T) {
				t.Parallel()

				var i2o wcnv.Converter = wcnv.ConvertFn(
					func(_ context.Context, jarr []byte) (output []byte, e error) {
						var strs []string
						e = json.Unmarshal(jarr, &strs)
						var joined string = strings.Join(strs, ",")
						return []byte(joined), e
					},
				).AsIf()

				type Reqs struct{ reqs []*rhp.Request }

				var ser2buf cser.ConvertFn[*Reqs] = func(
					serialized []byte,
					buf *Reqs,
				) error {
					var strs string = string(serialized)
					splited := strings.Split(strs, ",")
					buf.reqs = ua.Fold(
						splited,
						buf.reqs[:0],
						func(state []*rhp.Request, next string) []*rhp.Request {
							return append(state, &rhp.Request{Url: next})
						},
					)
					return nil
				}

				var msg2reqs cser.Message2requests[*Reqs] = func(
					ctx context.Context,
					msg *Reqs,
					dst chan<- pair.Pair[error, *rhp.Request],
				) error {
					return ua.TryForEach(
						msg.reqs,
						func(q *rhp.Request) error {
							select {
							case <-ctx.Done():
								return ctx.Err()
							default:
							}
							dst <- pair.Right[error](q)
							return nil
						},
					)
				}

				i2r := wser.Input2requests[*Reqs]{
					In2out:    i2o,
					Bytes2msg: ser2buf,
					Msg2reqs:  msg2reqs,
				}

				reqs := make(chan pair.Pair[error, *rhp.Request])

				go func() {
					defer close(reqs)
					var buf *Reqs = &Reqs{}

					e := i2r.Input2chan(
						context.Background(),
						[]byte(`["634","333","3776"]`),
						buf,
						reqs,
					)
					t.Run("no err", assertNil(e))
				}()

				var strs pair.Pair[error, []string] = uch.TryFold(
					context.Background(),
					nil,
					reqs,
					func(
						state []string,
						req *rhp.Request,
					) pair.Pair[error, []string] {
						return pair.Right[error](append(state, req.GetUrl()))
					},
				)

				t.Run("no chan err", assertNil(strs.Left))
				t.Run("2 items", assertEqual(len(strs.Right), 3))
				t.Run("tree", assertEqual(strs.Right[0], "634"))
				t.Run("tower", assertEqual(strs.Right[1], "333"))
				t.Run("fuji", assertEqual(strs.Right[2], "3776"))

			})
		})
	})
}
