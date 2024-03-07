package ch_test

import (
	"context"
	"testing"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
	cut "github.com/takanoriyanagitani/go-reqs2http/util/ch"
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

func TestChUtil(t *testing.T) {
	t.Parallel()

	t.Run("WriteToCh", func(t *testing.T) {
		t.Parallel()

		t.Run("ToChan", func(t *testing.T) {
			t.Parallel()

			t.Run("WriteToChFromSlice", func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					var empty []int32
					var w2c cut.WriteToCh[int32] = cut.WriteToChFromSlice(empty)
					var ctx context.Context = context.Background()
					var c <-chan pair.Pair[error, int32] = w2c.ToChan(ctx)
					for pair := range c {
						t.Errorf("not empty. got: %v\n", pair)
					}
				})
			})
		})

		t.Run("Fold", func(t *testing.T) {
			t.Parallel()

			t.Run("WriteToChFromSlice", func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					var empty []int32
					var w2c cut.WriteToCh[int32] = cut.WriteToChFromSlice(empty)
					var ctx context.Context = context.Background()
					var res pair.Pair[error, int32] = w2c.TryFold(
						ctx,
						0,
						func(state int32, next int32) pair.Pair[error, int32] {
							return pair.Right[error](state + next)
						},
					)
					t.Run("no error", assertNil(res.Left))
					t.Run("no items", assertEqual(res.Right, 0))
				})

				t.Run("integers", func(t *testing.T) {
					t.Parallel()

					var integers []uint32 = []uint32{
						123000,
						456,
					}
					var w2c cut.WriteToCh[uint32] = cut.WriteToChFromSlice(integers)
					var ctx context.Context = context.Background()
					var res pair.Pair[error, uint32] = w2c.TryFold(
						ctx,
						0,
						func(state uint32, next uint32) pair.Pair[error, uint32] {
							return pair.Right[error](state + next)
						},
					)
					t.Run("no error", assertNil(res.Left))
					t.Run("same tot", assertEqual(res.Right, 123456))
				})
			})
		})
	})
}
