package wser_test

import (
	"testing"
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
