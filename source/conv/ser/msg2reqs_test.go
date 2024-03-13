package cser_test

import (
	"context"
	"testing"

	"slices"
	"strings"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"
	ua "github.com/takanoriyanagitani/go-reqs2http/util/arr"
	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"

	cser "github.com/takanoriyanagitani/go-reqs2http/source/conv/ser"
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

func assertEqualItems[T comparable](a, b []T) func(*testing.T) {
	return assertEqualNew(func(a, b []T) (same bool) {
		return slices.Equal(a, b)
	})(a, b)
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

func TestMsg2reqs(t *testing.T) {
	t.Parallel()

	t.Run("Message2requests", func(t *testing.T) {
		t.Parallel()

		t.Run("Bytes2Chan", func(t *testing.T) {
			t.Parallel()

			t.Run("strings2reqs", func(t *testing.T) {
				t.Parallel()

				type Urls struct{ urls []string }

				var s2r cser.Message2requests[*Urls] = func(
					ctx context.Context,
					msg *Urls,
					dst chan<- pair.Pair[error, *rhp.Request],
				) error {
					return ua.TryForEach(
						msg.urls,
						func(url string) error {
							select {
							case <-ctx.Done():
								return ctx.Err()
							default:
							}

							dst <- pair.Right[error](&rhp.Request{Url: url})
							return nil
						},
					)
				}

				var s2u cser.ConvertFn[*Urls] = func(
					serialized []byte,
					urls *Urls,
				) error {
					const baseURL string = "https://localhost/"
					var s string = string(serialized)
					var splited []string = strings.Split(s, ",")
					var filtered = slices.DeleteFunc(
						splited,
						func(u string) (rmv bool) { return len(u) < 1 },
					)
					urls.urls = urls.urls[:0]
					return ua.TryForEach(
						filtered,
						func(url string) error {
							var u string = baseURL + url
							urls.urls = append(urls.urls, u)
							return nil
						},
					)
				}

				(func(m2r cser.Message2requests[*Urls]) {
					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						var buf Urls

						dst := make(chan pair.Pair[error, *rhp.Request])

						go func() {
							defer close(dst)
							e := m2r.Bytes2Chan(
								context.Background(),
								[]byte(""),
								s2u,
								&buf,
								dst,
							)
							if nil != e {
								panic(e)
							}
						}()

						var reqs pair.Pair[error, []*rhp.Request] = uch.TryFold(
							context.Background(),
							nil,
							dst,
							func(
								state []*rhp.Request,
								next *rhp.Request,
							) pair.Pair[error, []*rhp.Request] {
								return pair.Right[error](append(state, next))
							},
						)

						t.Run("no err", assertNil(reqs.Left))
						t.Run("no items", assertEqual(len(reqs.Right), 0))
					})

					t.Run("few", func(t *testing.T) {
						t.Parallel()

						var buf Urls

						dst := make(chan pair.Pair[error, *rhp.Request])

						go func() {
							defer close(dst)
							e := m2r.Bytes2Chan(
								context.Background(),
								[]byte("1,2,3"),
								s2u,
								&buf,
								dst,
							)
							if nil != e {
								panic(e)
							}
						}()

						var reqs pair.Pair[error, []*rhp.Request] = uch.TryFold(
							context.Background(),
							nil,
							dst,
							func(
								state []*rhp.Request,
								next *rhp.Request,
							) pair.Pair[error, []*rhp.Request] {
								return pair.Right[error](append(state, next))
							},
						)

						t.Run("no err", assertNil(reqs.Left))
						t.Run("3 items", assertEqual(len(reqs.Right), 3))

						t.Run("same items", assertEqualItems(
							ua.Map(
								reqs.Right,
								func(r *rhp.Request) string { return r.Url },
							),
							[]string{
								"https://localhost/1",
								"https://localhost/2",
								"https://localhost/3",
							},
						))
					})
				})(s2r)
			})
		})
	})
}
