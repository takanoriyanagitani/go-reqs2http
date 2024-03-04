package sendmany_test

import (
	"context"
	"testing"
	"time"

	buf "github.com/takanoriyanagitani/go-reqs2http/buffered"
	src "github.com/takanoriyanagitani/go-reqs2http/source"

	aps "github.com/takanoriyanagitani/go-reqs2http/app/sendmany"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

func must[T any](t T, e error) T {
	if nil == e {
		return t
	}

	panic(e)
}

func assertEqualNew[T any](
	comp func(a, b T) (same bool),
) func(a, b T) func(*testing.T) {
	return func(a, b T) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()

			var same bool = comp(a, b)
			switch same {
			case true:
				return
			default:
				t.Errorf("unexpected value got.\n")
				t.Errorf("expected: %v", b)
				t.Fatalf("got:      %v", a)
			}
		}
	}
}

func assertEqual[T comparable](a, b T) func(*testing.T) {
	return assertEqualNew(func(a, b T) (same bool) { return a == b })(a, b)
}

func assertNil(e error) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if nil == e {
			return
		}

		t.Fatalf("unexpected error: %v\n", e)
	}
}

func TestSender(t *testing.T) {
	t.Parallel()

	t.Run("ManySender", func(t *testing.T) {
		t.Parallel()

		t.Run("ManySenderNopDefault", func(t *testing.T) {
			t.Parallel()

			var ms aps.ManySender = aps.ManySenderNopDefault

			t.Run("nop", func(t *testing.T) {
				t.Parallel()

				var ctx context.Context = context.Background()
				e := ms.SendAll(ctx, 0)
				t.Run("no error", assertNil(e))
			})

			t.Run("short wait", func(t *testing.T) {
				t.Parallel()

				const wait time.Duration = 1 * time.Millisecond
				var waiter buf.WaitHint = buf.WaitHintFnStaticNew(wait)

				var wapp aps.ManySender = ms.WithWaiter(waiter)

				t.Run("count sender", func(t *testing.T) {
					t.Parallel()

					run := func(
						countUp func(),
						getCount func() int,
						reqs []*rhp.Request,
						ctx context.Context,
						expected int,
						noErr bool,
					) func(*testing.T) {
						return func(t *testing.T) {
							var counter buf.BufFn = func(
								_ context.Context,
								_ *rhp.Request,
							) error {
								countUp()
								return nil
							}

							var csend = buf.SenderFn{
								UsageFn:  buf.UsageSrcFnUnknown,
								ChangeFn: buf.ChangeSrcFnUnknown,
								BufferFn: counter,
							}

							var sapp aps.ManySender = wapp.WithSender(csend)

							var rsrc src.RequestSrcFn = src.
								RequestSrcFnFromSlice(reqs)

							var app aps.ManySender = sapp.WithSource(rsrc)
							e := app.SendAll(ctx, 0)

							t.Run("no error", assertEqual(noErr, nil == e))
							t.Run(
								"count check",
								assertEqual(getCount(), expected),
							)
						}
					}

					var dcnt int
					t.Run("double requests", run(
						func() { dcnt++ },
						func() int { return dcnt },
						[]*rhp.Request{
							{},
							{},
						},
						context.Background(),
						2,
						true,
					))

					var ccnt int
					cctx, ccan := context.WithCancel(context.Background())
					ccan()
					t.Run("cancel", run(
						func() { ccnt++ },
						func() int { return ccnt },
						[]*rhp.Request{
							{},
							{},
							{},
						},
						cctx,
						0,
						false,
					))

				})
			})
		})
	})
}
