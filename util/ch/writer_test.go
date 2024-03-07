package ch_test

import (
	"context"
	"errors"
	"testing"

	"encoding/json"
	"strconv"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
	util "github.com/takanoriyanagitani/go-reqs2http/util"
	ua "github.com/takanoriyanagitani/go-reqs2http/util/arr"

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

func assertErr(e error) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if nil != e {
			return
		}

		t.Fatalf("no error\n")
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

		t.Run("TryFold", func(t *testing.T) {
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

			t.Run("WriteToChFromBuilder", func(t *testing.T) {
				t.Parallel()

				t.Run("strings2integers", func(t *testing.T) {
					t.Parallel()

					t.Run("no err", func(t *testing.T) {
						t.Parallel()

						builder := func(
							_ context.Context,
						) (cut.WriteToCh[int8], error) {
							integers := []int8{42, -7}
							return cut.WriteToChFromSlice(integers), nil
						}

						w2c := cut.WriteToChFromBuilder(builder)
						var ctx context.Context = context.Background()
						var res pair.Pair[error, int8] = w2c.TryFold(
							ctx,
							0,
							func(state int8, next int8) pair.Pair[error, int8] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("no error", assertNil(res.Left))
						t.Run("same tot", assertEqual(res.Right, 35))
					})

					t.Run("err", func(t *testing.T) {
						t.Parallel()

						builder := func(
							_ context.Context,
						) (cut.WriteToCh[int8], error) {
							return nil, errors.New("invalid string")
						}

						w2c := cut.WriteToChFromBuilder(builder)
						var ctx context.Context = context.Background()
						var res pair.Pair[error, int8] = w2c.TryFold(
							ctx,
							0,
							func(state int8, next int8) pair.Pair[error, int8] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("error", assertErr(res.Left))
					})
				})
			})

			t.Run("Write2ChConverter", func(t *testing.T) {
				t.Parallel()

				t.Run("json2integers", func(t *testing.T) {
					t.Parallel()

					var w2cc cut.Write2ChConverter[string, int8] = func(
						_ context.Context,
						jstr string,
					) (cut.WriteToCh[int8], error) {
						var parsed []int8
						e := json.Unmarshal([]byte(jstr), &parsed)
						if nil != e {
							return nil, e
						}
						var w2c cut.WriteToCh[int8] = cut.WriteToChFromSlice(
							parsed,
						)
						return w2c, nil
					}

					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						var ctx context.Context = context.Background()
						w2c := w2cc.ToConverted(`[]`)
						var res pair.Pair[error, int8] = w2c.TryFold(
							ctx,
							0,
							func(state int8, next int8) pair.Pair[error, int8] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("no err", assertNil(res.Left))
						t.Run("no items", assertEqual(res.Right, 0))
					})

					t.Run("integers", func(t *testing.T) {
						t.Parallel()

						var ctx context.Context = context.Background()
						w2c := w2cc.ToConverted(`[7, -42]`)
						var res pair.Pair[error, int8] = w2c.TryFold(
							ctx,
							0,
							func(state int8, next int8) pair.Pair[error, int8] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("no err", assertNil(res.Left))
						t.Run("no items", assertEqual(res.Right, -35))
					})

					t.Run("invalid json", func(t *testing.T) {
						t.Parallel()

						var ctx context.Context = context.Background()
						w2c := w2cc.ToConverted(`]`)
						var res pair.Pair[error, int8] = w2c.TryFold(
							ctx,
							0,
							func(state int8, next int8) pair.Pair[error, int8] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("err", assertErr(res.Left))
					})
				})

				t.Run("json2strs2ints", func(t *testing.T) {
					t.Parallel()

					var w2cc cut.Write2ChConverter[[]string, int] = func(
						_ context.Context,
						strs []string,
					) (cut.WriteToCh[int], error) {
						parsed, e := ua.MapErr(
							strs,
							strconv.Atoi,
						)
						w2c := cut.WriteToChFromSlice(parsed)
						return w2c, e
					}

					json2strings := func(jstr string) ([]string, error) {
						var ret []string
						e := json.Unmarshal([]byte(jstr), &ret)
						return ret, e
					}

					var json2w2c func(
						ctx context.Context,
						jstr string,
					) (cut.WriteToCh[int], error) = util.ComposeCtx(
						util.CtxIgnore(json2strings),
						w2cc,
					)

					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						var ctx context.Context = context.Background()
						w2c, e := json2w2c(ctx, `[]`)
						t.Run("no err", assertNil(e))
						var res pair.Pair[error, int] = w2c.TryFold(
							ctx,
							0,
							func(state int, next int) pair.Pair[error, int] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("no err", assertNil(res.Left))
						t.Run("no items", assertEqual(res.Right, 0))
					})

					t.Run("integers", func(t *testing.T) {
						t.Parallel()

						var ctx context.Context = context.Background()
						w2c, e := json2w2c(ctx, `["42", "-20"]`)
						t.Run("no err", assertNil(e))
						var res pair.Pair[error, int] = w2c.TryFold(
							ctx,
							0,
							func(state int, next int) pair.Pair[error, int] {
								return pair.Right[error](state + next)
							},
						)
						t.Run("no err", assertNil(res.Left))
						t.Run("same value", assertEqual(res.Right, 22))
					})
				})
			})
		})
	})
}
