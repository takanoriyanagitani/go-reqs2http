package reqs2http_test

import (
	"net/http"
	"testing"

	r2h "github.com/takanoriyanagitani/go-reqs2http"
	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

func assertEqualNew[T any](comp func(a, b T) (same bool)) func(a, b T) func(*testing.T) {
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

func TestMethod(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		table := []struct {
			Method   rhp.Method
			Expected string
		}{
			{rhp.Method_METHOD_UNSPECIFIED, http.MethodGet},
			{rhp.Method_METHOD_GET, http.MethodGet},
			{rhp.Method_METHOD_HEAD, http.MethodHead},
			{rhp.Method_METHOD_POST, http.MethodPost},
			{rhp.Method_METHOD_PUT, http.MethodPut},
			{rhp.Method_METHOD_PATCH, http.MethodPatch},
			{rhp.Method_METHOD_DELETE, http.MethodDelete},
			{rhp.Method_METHOD_CONNECT, http.MethodConnect},
			{rhp.Method_METHOD_OPTIONS, http.MethodOptions},
			{rhp.Method_METHOD_TRACE, http.MethodTrace},
		}

		for _, tab := range table {
			var method rhp.Method = tab.Method
			var expected string = tab.Expected
			m := r2h.Method{method}
			var s string = m.String()
			t.Run("same string", assertEqual(s, expected))
		}
	})
}
