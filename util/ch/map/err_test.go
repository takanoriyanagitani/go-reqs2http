package chmap_test

import (
	"context"
	"testing"

	"strconv"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
	util "github.com/takanoriyanagitani/go-reqs2http/util"

	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"
	chm "github.com/takanoriyanagitani/go-reqs2http/util/ch/map"
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

func assertFalse(b bool) func(*testing.T) {
	return assertEqual(b, false)
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

func assertErr(e error) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if nil != e {
			return
		}

		t.Fatalf("no error\n")
	}
}

func TestErr(t *testing.T) {
	t.Parallel()

	t.Run("MapErr", func(t *testing.T) {
		t.Parallel()

		t.Run("strings2integers", func(t *testing.T) {
			t.Parallel()

			var mapper func(context.Context, string) (int, error) = util.
				CtxIgnore(strconv.Atoi)

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				strings := make(chan string)
				close(strings)

				var pairs <-chan pair.Pair[error, string] = chm.Map(
					context.Background(),
					strings,
					func(_ context.Context, s string) pair.Pair[error, string] {
						return pair.Right[error](s)
					},
				)

				var mapd <-chan pair.Pair[error, int] = chm.MapErr(
					context.Background(),
					pairs,
					mapper,
				)

				_, ok := <-mapd
				t.Run("no items", assertFalse(ok))
			})

			t.Run("integers", func(t *testing.T) {
				t.Parallel()

				strings := make(chan string)
				go func() {
					defer close(strings)

					strings <- "1"
					strings <- "42"
				}()

				var pairs <-chan pair.Pair[error, string] = chm.Map(
					context.Background(),
					strings,
					func(_ context.Context, s string) pair.Pair[error, string] {
						return pair.Right[error](s)
					},
				)

				var mapd <-chan pair.Pair[error, int] = chm.MapErr(
					context.Background(),
					pairs,
					mapper,
				)

				var res pair.Pair[error, int] = uch.TryFold(
					context.Background(),
					0,
					mapd,
					func(state int, next int) pair.Pair[error, int] {
						return pair.Right[error](state + next)
					},
				)
				t.Run("no error", assertNil(res.Left))
				t.Run("same value", assertEqual(res.Right, 43))
			})
		})
	})

	t.Run("ConvErr", func(t *testing.T) {
		t.Parallel()

		t.Run("strings2ints2floats", func(t *testing.T) {
			t.Parallel()

			var s2i func(context.Context, string) (int, error) = util.
				CtxIgnore(strconv.Atoi)
			i2f := func(_ context.Context, i int) (float64, error) {
				return float64(i), nil
			}

			var s2f chm.ConvErr[string, float64] = util.ComposeCtx(
				s2i,
				i2f,
			)

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				strs := make(chan string)
				close(strs)
				var pairs <-chan pair.Pair[error, string] = chm.Map(
					context.Background(),
					strs,
					func(_ context.Context, s string) pair.Pair[error, string] {
						return pair.Right[error](s)
					},
				)

				var mapd <-chan pair.Pair[error, float64] = s2f.MapErr(
					context.Background(),
					pairs,
				)

				var res pair.Pair[error, float64] = uch.TryFold(
					context.Background(),
					0.0,
					mapd,
					func(
						state float64,
						next float64,
					) pair.Pair[error, float64] {
						return pair.Right[error](state + next)
					},
				)
				t.Run("no err", assertNil(res.Left))
				t.Run("no items", assertEqual(res.Right, 0.0))
			})

			t.Run("strings", func(t *testing.T) {
				t.Parallel()

				strs := make(chan string)
				go func() {
					defer close(strs)

					strs <- "3776"
					strs <- "599"
				}()
				var pairs <-chan pair.Pair[error, string] = chm.Map(
					context.Background(),
					strs,
					func(_ context.Context, s string) pair.Pair[error, string] {
						return pair.Right[error](s)
					},
				)

				var mapd <-chan pair.Pair[error, float64] = s2f.MapErr(
					context.Background(),
					pairs,
				)

				var res pair.Pair[error, float64] = uch.TryFold(
					context.Background(),
					0.0,
					mapd,
					func(
						state float64,
						next float64,
					) pair.Pair[error, float64] {
						return pair.Right[error](state + next)
					},
				)
				t.Run("no err", assertNil(res.Left))
				t.Run("same value", assertEqual(res.Right, 4375.0))
			})
		})
	})
}
