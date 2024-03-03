package source_test

import (
	"errors"
	"testing"

	"context"

	src "github.com/takanoriyanagitani/go-reqs2http/source"

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

func assertEmpty[T any](s []T) func(*testing.T) {
	return assertEqual(len(s), 0)
}

func assertTrue(b bool) func(*testing.T) {
	return assertEqual(b, true)
}

func TestReqs(t *testing.T) {
	t.Parallel()

	t.Run("RequestSrcFn", func(t *testing.T) {
		t.Parallel()

		t.Run("RequestSrcFnEmpty", func(t *testing.T) {
			t.Parallel()

			var rsf src.RequestSrcFn = src.RequestSrcFnEmpty
			_, e := rsf(context.Background())
			t.Run("expected error", assertTrue(errors.Is(e, src.ErrNoMoreData)))
		})

		t.Run("RequestSrcFnFromSlice", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var rsf src.RequestSrcFn = src.RequestSrcFnFromSlice(nil)
				var rsc src.RequestSourceCh = rsf.ToChan(0)
				var rscf src.ReqSrcChanFn = rsc.GetRequests
				var ctx context.Context = context.Background()
				var pairs []src.RequestResult = rscf.GetAll(ctx)
				t.Run("empty slice", assertEmpty(pairs))
			})

			t.Run("single", func(t *testing.T) {
				t.Parallel()

				var rsf src.RequestSrcFn = src.RequestSrcFnFromSlice(
					[]*rhp.Request{{}},
				)
				var rsc src.RequestSourceCh = rsf.ToChan(0)
				var rscf src.ReqSrcChanFn = rsc.GetRequests
				var ctx context.Context = context.Background()
				var pairs []src.RequestResult = rscf.GetAll(ctx)
				t.Run("single item", assertEqual(len(pairs), 1))
			})

			t.Run("double", func(t *testing.T) {
				t.Parallel()

				var rsf src.RequestSrcFn = src.RequestSrcFnFromSlice(
					[]*rhp.Request{
						{},
						{},
					},
				)
				var rsc src.RequestSourceCh = rsf.ToChan(0)
				var rscf src.ReqSrcChanFn = rsc.GetRequests
				var ctx context.Context = context.Background()
				var pairs []src.RequestResult = rscf.GetAll(ctx)
				t.Run("double items", assertEqual(len(pairs), 2))
			})

			t.Run("cancel", func(t *testing.T) {
				t.Parallel()

				var rsf src.RequestSrcFn = src.RequestSrcFnFromSlice(
					[]*rhp.Request{
						{},
						{},
					},
				)
				var rsc src.RequestSourceCh = rsf.ToChan(0)
				var rscf src.ReqSrcChanFn = rsc.GetRequests
				var ctx context.Context = context.Background()
				ctx, can := context.WithCancel(ctx)
				can()
				var pairs []src.RequestResult = rscf.GetAll(ctx)
				t.Run("no items", assertEqual(len(pairs), 0))
			})
		})
	})

	t.Run("ReqSrcChanFn", func(t *testing.T) {
		t.Parallel()

		t.Run("GetAll", func(t *testing.T) {
			t.Parallel()

			t.Run("RequestSrcFnEmpty", func(t *testing.T) {
				t.Parallel()

				var rsf src.RequestSrcFn = src.RequestSrcFnEmpty
				var rsc src.RequestSourceCh = rsf.ToChan(0)
				var rscf src.ReqSrcChanFn = rsc.GetRequests
				var ctx context.Context = context.Background()
				var pairs []src.RequestResult = rscf.GetAll(ctx)
				t.Run("empty slice", assertEmpty(pairs))
			})

			t.Run("ReqSrcChanFnFromSlice", func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					var rscf src.ReqSrcChanFn = src.ReqSrcChanFnFromSlice(nil)
					var ctx context.Context = context.Background()
					var pairs []src.RequestResult = rscf.GetAll(ctx)
					t.Run("empty", assertEmpty(pairs))
				})

				t.Run("single", func(t *testing.T) {
					t.Parallel()

					var rscf src.ReqSrcChanFn = src.ReqSrcChanFnFromSlice(
						[]*rhp.Request{{}},
					)
					var ctx context.Context = context.Background()
					var pairs []src.RequestResult = rscf.GetAll(ctx)
					t.Run("single item", assertEqual(len(pairs), 1))
				})
			})
		})
	})
}
