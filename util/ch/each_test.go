package ch_test

import (
	"context"
	"testing"

	pair "github.com/takanoriyanagitani/go-reqs2http/pair"

	uch "github.com/takanoriyanagitani/go-reqs2http/util/ch"
)

func TestEach(t *testing.T) {
	t.Parallel()

	t.Run("TryForEach", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			pairs := make(chan pair.Pair[error, int])
			close(pairs)

			e := uch.TryForEach(
				context.Background(),
				pairs,
				func(_ int) error { return nil },
			)
			t.Run("no error", assertNil(e))
		})

		t.Run("done", func(t *testing.T) {
			t.Parallel()

			pairs := make(chan pair.Pair[error, int])
			defer close(pairs)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			e := uch.TryForEach(
				ctx,
				pairs,
				func(_ int) error { return nil },
			)
			t.Run("done", assertErr(e))
		})
	})
}
